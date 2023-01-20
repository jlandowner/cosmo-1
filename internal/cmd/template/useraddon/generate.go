package useraddon

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"

	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	cosmov1alpha1 "github.com/cosmo-workspace/cosmo/api/v1alpha1"
	"github.com/cosmo-workspace/cosmo/pkg/clog"
	cmdutil "github.com/cosmo-workspace/cosmo/pkg/cmdutil"
	"github.com/cosmo-workspace/cosmo/pkg/template"
)

type GenerateOption struct {
	*cmdutil.CliOptions

	Name         string
	OutputFile   string
	RequiredVars string
	Desc         string

	SetDefaultUserAddon bool
	DisableNamePrefix   bool

	ClusterScope bool

	tmpl  cosmov1alpha1.TemplateObject
	input []byte
}

func GenerateCmd(cmd *cobra.Command, cliOpt *cmdutil.CliOptions) *cobra.Command {
	o := &GenerateOption{CliOptions: cliOpt}
	cmd.PersistentPreRunE = o.PreRunE
	cmd.RunE = cmdutil.RunEHandler(o.RunE)

	cmd.Flags().StringVarP(&o.Name, "name", "n", "", "template name (use directory name if not specified)")
	cmd.Flags().StringVarP(&o.OutputFile, "output", "o", "", "write output into file (default: Stdout)")
	cmd.Flags().StringVar(&o.RequiredVars, "required-vars", "", "template custom vars to be replaced by instance. format --required-vars VAR1,VAR2:default-value")
	cmd.Flags().StringVar(&o.Desc, "desc", "", "template description")

	cmd.Flags().BoolVar(&o.SetDefaultUserAddon, "set-default-user-addon", false, "set default user addon")
	cmd.Flags().BoolVar(&o.DisableNamePrefix, "disable-nameprefix", false, "disable adding instance name prefix on child resource name")
	cmd.Flags().BoolVar(&o.ClusterScope, "cluster-scope", false, "generate ClusterTemplate (default generate namespaced Template)")

	return cmd
}

func (o *GenerateOption) PreRunE(cmd *cobra.Command, args []string) error {
	if err := o.Validate(cmd, args); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}
	if err := o.Complete(cmd, args); err != nil {
		return fmt.Errorf("invalid options: %w", err)
	}
	return nil
}

func (o *GenerateOption) Validate(cmd *cobra.Command, args []string) error {
	if err := o.CliOptions.Validate(cmd, args); err != nil {
		return err
	}

	if isatty.IsTerminal(os.Stdin.Fd()) {
		return fmt.Errorf("no input via stdin")
	}

	return nil
}

func (o *GenerateOption) Complete(cmd *cobra.Command, args []string) error {
	o.CliOptions.CompleteLogger()

	// Complete ClusterTemplate or not
	if o.ClusterScope {
		o.tmpl = &cosmov1alpha1.ClusterTemplate{}
	} else {
		o.tmpl = &cosmov1alpha1.Template{}
	}

	// Complete Name from direcotry name is not specified
	if o.Name == "" {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		o.Name = filepath.Base(dir)
	}

	// Complete OutputFile path if specified
	if o.OutputFile != "" {
		var err error
		o.OutputFile, err = filepath.Abs(o.OutputFile)
		if err != nil {
			return err
		}
	}

	// Format RequiredVars if specified
	if o.RequiredVars != "" {
		varsList := strings.Split(o.RequiredVars, ",")

		vars := make([]cosmov1alpha1.RequiredVarSpec, 0, len(varsList))
		for _, v := range varsList {
			vcol := strings.Split(v, ":")
			varSpec := cosmov1alpha1.RequiredVarSpec{Var: vcol[0]}
			if len(vcol) > 1 {
				varSpec.Default = vcol[1]
			}
			vars = append(vars, varSpec)
		}
		o.tmpl.GetSpec().RequiredVars = vars
	}

	// Set metadata
	o.tmpl.SetName(o.Name)
	o.tmpl.GetSpec().Description = o.Desc

	// Set GroupVersionKind for generating YAML
	scheme := runtime.NewScheme()
	cosmov1alpha1.AddToScheme(scheme)
	gvk, err := apiutil.GVKForObject(o.tmpl, scheme)
	if err != nil {
		return err
	}
	o.tmpl.SetGroupVersionKind(gvk)

	// Set annotation for specific type
	template.SetTemplateType(o.tmpl, cosmov1alpha1.TemplateLabelEnumTypeUserAddon)

	ann := o.tmpl.GetAnnotations()
	if ann == nil {
		ann = make(map[string]string)
	}
	if o.SetDefaultUserAddon {
		ann[cosmov1alpha1.UserAddonTemplateAnnKeyDefaultUserAddon] = strconv.FormatBool(true)
	}
	if o.DisableNamePrefix {
		ann[cosmov1alpha1.TemplateAnnKeyDisableNamePrefix] = strconv.FormatBool(true)
	}
	o.tmpl.SetAnnotations(ann)

	// input data from stdin
	o.input, err = io.ReadAll(o.In)
	if err != nil {
		return fmt.Errorf("failed to read input file : %w", err)
	}
	if len(o.input) == 0 {
		return fmt.Errorf("no input")
	}

	return nil
}

func (o *GenerateOption) RunE(cmd *cobra.Command, args []string) error {
	log := o.Logr.WithName("gen_useraddon_template")
	ctx := clog.IntoContext(o.Ctx, log)

	// create tmp dir
	tmpDir, err := ioutil.TempDir(os.TempDir(), "cosmoctl-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir : %w", err)
	}
	defer os.RemoveAll(tmpDir)
	log.Debug().Info("tmpDir created", "path", tmpDir)

	kust := cmdutil.NewKustomize(tmpDir)
	if o.DisableNamePrefix {
		kust.WithDisableNamePrefix()
	}

	if _, err := kust.WritePackagedFile(o.input); err != nil {
		return err
	}
	log.Debug().Info("packaged file created")

	// run kustomize
	out, err := kust.Build(ctx)
	if err != nil {
		return err
	}
	log.Debug().Info("executed kustomize", "output", string(out))

	o.tmpl.GetSpec().RawYaml = string(out)
	outtmpl, err := yaml.Marshal(&o.tmpl)
	if err != nil {
		return err
	}

	output := append([]byte("# Generated by cosmoctl template command\n"), outtmpl...)

	// output to Stdout or write the output to file given by option
	if o.OutputFile == "" {
		fmt.Fprintln(o.Out, string(output))
	} else {
		if err := cmdutil.CreateFile(filepath.Dir(o.OutputFile), filepath.Base(o.OutputFile), output); err != nil {
			return err
		}
	}
	return nil
}
