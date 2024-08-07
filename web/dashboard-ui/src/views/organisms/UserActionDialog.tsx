import {
  Add,
  Close,
  ExpandLess,
  ExpandMore,
  PersonOutlineTwoTone,
  SecurityOutlined,
} from "@mui/icons-material";
import {
  Alert,
  Box,
  Button,
  Checkbox,
  Chip,
  Collapse,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  FormControlLabel,
  FormHelperText,
  Grid,
  IconButton,
  InputAdornment,
  List,
  ListItem,
  ListItemText,
  MenuItem,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  TextField,
  Tooltip,
  Typography,
} from "@mui/material";
import React, { useEffect, useState } from "react";
import { UseFormRegisterReturn, useFieldArray, useForm } from "react-hook-form";
import { DialogContext } from "../../components/ContextProvider";
import { Template } from "../../proto/gen/dashboard/v1alpha1/template_pb";
import { User, UserAddon } from "../../proto/gen/dashboard/v1alpha1/user_pb";
import { NameAvatar } from "../atoms/NameAvatar";
import { FormSelectableChip } from "../atoms/SelectableChips";
import { TextFieldLabel } from "../atoms/TextFieldLabel";
import { PasswordDialogContext } from "./PasswordDialog";
import {
  isAdminRole,
  isPrivilegedRole,
  useTemplates,
  useUserModule,
} from "./UserModule";

const registerMui = ({ ref, ...rest }: UseFormRegisterReturn) => ({
  inputRef: ref,
  ...rest,
});

/**
 * UserActionDialog
 */
interface UserActionDialogProps {
  title: string;
  actions: React.ReactNode;
  user: User;
  onClose?: () => void;
  defaultOpenUserAddon?: boolean;
}

const UserActionDialog: React.FC<UserActionDialogProps> = ({
  title,
  actions,
  user,
  onClose,
  defaultOpenUserAddon,
}) => {
  console.log(user);
  const [openUserAddon, setOpenUserAddon] = useState<boolean>(
    defaultOpenUserAddon || false
  );

  const handleOpenUserAddonClick = () => {
    setOpenUserAddon(!openUserAddon);
  };

  return (
    <Dialog open={true} onClose={onClose} fullWidth maxWidth={"xs"}>
      <DialogTitle>
        {title}
        <IconButton
          sx={{
            position: "absolute",
            right: 8,
            top: 8,
            color: (theme) => theme.palette.grey[500],
          }}
          onClick={onClose}
        >
          <Close />
        </IconButton>
      </DialogTitle>
      <DialogContent>
        <Stack spacing={2}>
          <Stack alignItems="center">
            <NameAvatar
              name={user?.displayName}
              sx={{ width: 50, height: 50 }}
            />
          </Stack>
          <TextFieldLabel
            label="Name"
            fullWidth
            value={user.name}
            startAdornmentIcon={<PersonOutlineTwoTone />}
          />
          <TextFieldLabel
            label="Display Name"
            fullWidth
            value={user.displayName}
            startAdornmentIcon={<PersonOutlineTwoTone />}
          />
          <TextFieldLabel
            label="AuthType"
            fullWidth
            value={user.authType}
            startAdornmentIcon={<SecurityOutlined />}
          />
          <Typography color="text.secondary" display="block" variant="caption">
            Roles
          </Typography>
          <Stack alignItems="center">
            <Grid container justifyContent="center" sx={{ width: 300 }}>
              {user?.roles &&
                user.roles.map((v, i) => {
                  return (
                    <Grid item key={i}>
                      <Chip
                        color={
                          isPrivilegedRole(v)
                            ? "error"
                            : isAdminRole(v)
                            ? "warning"
                            : "default"
                        }
                        size="small"
                        key={i}
                        label={v}
                        sx={{ m: 0.05 }}
                      />
                    </Grid>
                  );
                })}
            </Grid>
          </Stack>
          <Divider />
          {Boolean(user.addons.length) && (
            <Stack spacing={1}>
              <Typography
                color="text.secondary"
                display="block"
                variant="caption"
              >
                User Addons
                <IconButton
                  size="small"
                  aria-label="openUserAddon"
                  onClick={handleOpenUserAddonClick}
                >
                  {openUserAddon ? (
                    <ExpandLess fontSize="small" />
                  ) : (
                    <ExpandMore fontSize="small" />
                  )}
                </IconButton>
              </Typography>
              <Collapse in={openUserAddon} timeout="auto" unmountOnExit>
                <List component="nav" dense={true}>
                  {user.addons.map((v, i) => (
                    <ListItem key={i}>
                      <ListItemText
                        disableTypography={true}
                        primary={
                          <Typography
                            color="text.secondary"
                            display="block"
                            variant="caption"
                          >
                            * {v.template}
                          </Typography>
                        }
                        secondary={
                          <TableContainer component={Paper}>
                            <Table aria-label={v.template}>
                              <TableBody>
                                {Object.keys(v.vars).map((key, j) => (
                                  <TableRow
                                    key={j}
                                    sx={{
                                      "&:last-child td, &:last-child th": {
                                        border: 0,
                                      },
                                    }}
                                  >
                                    <TableCell component="th" scope="row">
                                      {key}
                                    </TableCell>
                                    <TableCell align="right">
                                      {v.vars[key]}
                                    </TableCell>
                                  </TableRow>
                                ))}
                              </TableBody>
                            </Table>
                          </TableContainer>
                        }
                      />
                    </ListItem>
                  ))}
                </List>
              </Collapse>
            </Stack>
          )}
        </Stack>
      </DialogContent>
      <DialogActions>{actions}</DialogActions>
    </Dialog>
  );
};

/**
 * Delete
 */
export const UserDeleteDialog: React.VFC<{
  onClose: () => void;
  user: User;
}> = ({ onClose, user }) => {
  console.log("UserDeleteDialog");
  const hooks = useUserModule();
  const [lock, setLock] = useState(false);

  return (
    <UserActionDialog
      title="Delete User 👋"
      onClose={onClose}
      user={user}
      actions={
        <Alert
          severity="warning"
          action={
            <>
              <Checkbox
                color="warning"
                onChange={(e) => setLock(e.target.checked)}
              />
              <Button
                variant="contained"
                color="secondary"
                disabled={!lock}
                onClick={() => {
                  hooks.deleteUser(user.name).then(() => onClose());
                }}
              >
                Delete
              </Button>
            </>
          }
        >
          This action is NOT recoverable. Are you sure to delete it?
        </Alert>
      }
    />
  );
};

export const UserCreateConfirmDialog: React.VFC<{
  onClose: () => void;
  onConfirm: () => void;
  user: User;
}> = ({ onClose, onConfirm, user }) => {
  console.log("UserCreateConfirmDialog");

  const hooks = useUserModule();
  const passwordDialogDispatch = PasswordDialogContext.useDispatch();
  return (
    <UserActionDialog
      title="Create?"
      user={user}
      defaultOpenUserAddon={true}
      onClose={onClose}
      actions={
        <DialogActions>
          <Button onClick={onClose} color="primary">
            Back
          </Button>
          <Button
            variant="contained"
            color="secondary"
            onClick={() => {
              hooks
                .createUser(
                  user.name,
                  user.displayName,
                  user.authType,
                  user.roles,
                  user.addons
                )
                .then((newUser) => {
                  onClose();
                  onConfirm();
                  if (newUser?.defaultPassword) {
                    passwordDialogDispatch(true, { user: newUser! });
                  }
                  hooks.getUsers();
                });
            }}
          >
            Create
          </Button>
        </DialogActions>
      }
    />
  );
};

/**
 * Create
 */
type Inputs = {
  id: string;
  name: string;
  authType: string;
  roles: { name: string; enabled: boolean }[];
  customRole: string;
  addons: {
    template: Template;
    enable: boolean;
    vars: string[];
  }[];
};
export const UserCreateDialog: React.VFC<{ onClose: () => void }> = ({
  onClose,
}) => {
  console.log("UserCreateDialog");
  const hooks = useUserModule();
  const userCreateConfirmDialogDispatch =
    UserCreateConfirmDialogContext.useDispatch();

  const {
    register,
    handleSubmit,
    watch,
    control,
    formState: { errors },
    getValues,
    setValue,
    setError,
  } = useForm<Inputs>({
    defaultValues: {
      roles: hooks.existingRoles.map((v) => ({
        name: v,
        enabled: false,
      })),
      customRole: "",
    },
  });

  const { fields: addonsFields, replace: replaceAddons } = useFieldArray({
    control,
    name: "addons",
  });

  const templ = useTemplates();
  useEffect(() => {
    templ.getUserAddonTemplates({ useRoleFilter: true });
  }, []); // eslint-disable-line
  useEffect(() => {
    replaceAddons(
      templ.templates.map((t) => ({ template: t, enable: false, vars: [] }))
    );
  }, [templ.templates]); // eslint-disable-line

  const { fields: rolesFields, append: appendRoles } = useFieldArray({
    control,
    name: "roles",
    rules: {
      validate: (fieldArrayValues) => {
        // check that no duplicates exist
        const values = fieldArrayValues
          .map((item) => item.name)
          .filter((v) => v !== "");
        const uniqueValues = [...new Set(values)];
        return values.length === uniqueValues.length || "No duplicates allowed";
      },
    },
  });

  const [openCustomInput, setOpenCustomInput] = useState<boolean>(false);

  return (
    <Dialog open={true} fullWidth>
      <DialogTitle>Create New User 🎉</DialogTitle>
      <form
        onSubmit={handleSubmit((inp: Inputs) => {
          console.log(inp);

          const enabled = inp.addons.filter(
            (v) => v.enable || v.template.isDefaultUserAddon
          );

          // check required addons are enabled
          for (const hasDeps of enabled.filter(
            (v) => v.template.requiredUseraddons.length > 0
          )) {
            for (const req of hasDeps.template.requiredUseraddons) {
              if (enabled.findIndex((v) => v.template.name === req) < 0) {
                setError(
                  `addons.${inp.addons.findIndex(
                    (v) => v.template.name === hasDeps.template.name
                  )}.enable`,
                  {
                    message: `Required: ${hasDeps.template.requiredUseraddons.join(
                      ", "
                    )}`,
                  }
                );
                return;
              }
            }
          }

          const userAddons = enabled.map((inpAddon) => {
            const vars: { [key: string]: string } = {};
            inpAddon.vars.forEach((v, i) => {
              vars[inpAddon.template.requiredVars?.[i].varName!] = v;
            });
            return {
              template: inpAddon.template.name,
              vars: vars,
              clusterScoped: inpAddon.template.isClusterScope,
            };
          });
          const protoUserAddons = userAddons.map((ua) => new UserAddon(ua));
          console.log("protoUserAddons", protoUserAddons);

          let protoRoles = inp.roles
            .filter((v) => v.enabled)
            .map((v) => v.name);
          protoRoles = [...new Set(protoRoles)]; // remove duplicates
          console.log("protoRoles", protoRoles);

          userCreateConfirmDialogDispatch(true, {
            onConfirm: () => {
              onClose();
            },
            user: new User({
              name: inp.id,
              displayName: inp.name,
              roles: protoRoles,
              authType: inp.authType,
              addons: protoUserAddons,
            }),
          });
        })}
        autoComplete="new-password"
      >
        <DialogContent>
          <Stack spacing={2}>
            <TextField
              label="User ID"
              fullWidth
              autoFocus
              {...registerMui(
                register("id", {
                  required: { value: true, message: "Required" },
                  pattern: {
                    value: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/,
                    message:
                      'Only lowercase alphanumeric characters or "-" are allowed (start and end with an alphanumeric character)',
                  },
                  maxLength: { value: 50, message: "Max 50 characters" },
                })
              )}
              error={Boolean(errors.id)}
              helperText={
                (errors.id && errors.id.message) ||
                'Lowercase Alphanumeric or "-"'
              }
              InputProps={{
                autoComplete: "off",
                startAdornment: (
                  <InputAdornment position="start">
                    <PersonOutlineTwoTone />
                  </InputAdornment>
                ),
              }}
            />
            <TextField
              label="Display Name (optional)"
              fullWidth
              {...registerMui(
                register("name", {
                  maxLength: { value: 50, message: "Max 50 characters" },
                })
              )}
              error={Boolean(errors.name)}
              helperText={errors.name && errors.name.message}
              InputProps={{
                autoComplete: "off",
                startAdornment: (
                  <InputAdornment position="start">
                    <PersonOutlineTwoTone />
                  </InputAdornment>
                ),
              }}
            />
            <TextField
              label="Auth Type"
              select
              fullWidth
              defaultValue="password-secret"
              {...registerMui(
                register("authType", {
                  required: { value: true, message: "Required" },
                })
              )}
              error={Boolean(errors.authType)}
              helperText={errors.authType && errors.authType.message}
              InputProps={{
                autoComplete: "off",
                startAdornment: (
                  <InputAdornment position="start">
                    <SecurityOutlined />
                  </InputAdornment>
                ),
              }}
            >
              <MenuItem key={"password-secret"} value={"password-secret"}>
                <Tooltip
                  title={"Authentication by password registered with cosmo"}
                  placement="right"
                  arrow
                  enterDelay={500}
                >
                  <div>password-secret</div>
                </Tooltip>
              </MenuItem>
              <MenuItem key={"ldap"} value={"ldap"}>
                <Tooltip
                  title={"Authentication by ldap"}
                  placement="right"
                  arrow
                  enterDelay={500}
                >
                  <div>ldap</div>
                </Tooltip>
              </MenuItem>
            </TextField>
            <Typography
              color="text.secondary"
              display="block"
              variant="caption"
            >
              Roles
            </Typography>
            <Grid container>
              {rolesFields.map((v, index) => (
                <FormSelectableChip
                  defaultChecked={v.enabled}
                  key={index}
                  control={control}
                  label={v.name}
                  color="primary"
                  sx={{ m: 0.05 }}
                  {...register(`roles.${index}.enabled` as const)}
                />
              ))}
            </Grid>
            <FormHelperText error={Boolean(errors.roles?.root?.message)}>
              {errors.roles?.root?.message}
            </FormHelperText>
            <Box display="flex" alignItems="center">
              <Typography color="text.secondary" variant="caption">
                Add Custom Role
              </Typography>
              <IconButton
                size="small"
                onClick={() => setOpenCustomInput(!openCustomInput)}
              >
                {openCustomInput ? <ExpandLess /> : <Add />}
              </IconButton>
            </Box>
            {openCustomInput && (
              <TextField
                label="Custom Role"
                {...register(`customRole`)}
                InputProps={{
                  endAdornment: (
                    <IconButton
                      size="small"
                      onClick={() => {
                        if (!getValues(`customRole`)) return;
                        if (
                          rolesFields
                            .map((v) => v.name)
                            .includes(getValues(`customRole`))
                        ) {
                          setError(`customRole`, {
                            message: "Role already exists",
                          });
                          return;
                        }
                        appendRoles({
                          name: getValues(`customRole`),
                          enabled: true,
                        });
                        setValue(`customRole`, "");
                        setError(`customRole`, {});
                      }}
                    >
                      <Add />
                    </IconButton>
                  ),
                }}
              />
            )}
            {Boolean(errors.customRole?.message) && (
              <FormHelperText error={Boolean(errors.customRole?.message)}>
                {errors.customRole?.message}
              </FormHelperText>
            )}
            <Divider />
            <Stack spacing={1}>
              {Boolean(templ.templates.length) && (
                <Typography
                  color="text.secondary"
                  display="block"
                  variant="caption"
                >
                  Enable User Addons
                </Typography>
              )}
              {addonsFields.map((field, index) => (
                <React.Fragment key={field.id}>
                  <FormControlLabel
                    label={field.template.name}
                    control={
                      <Tooltip
                        title={field.template.description || "No description"}
                        placement="bottom"
                        arrow
                        enterDelay={1000}
                      >
                        <Checkbox
                          defaultChecked={
                            field.template.isDefaultUserAddon || false
                          }
                          disabled={field.template.isDefaultUserAddon || false}
                          {...registerMui(
                            register(`addons.${index}.enable` as const, {})
                          )}
                        />
                      </Tooltip>
                    }
                  />
                  <FormHelperText
                    error={Boolean(errors.addons?.[index]?.enable)}
                  >
                    {errors.addons?.[index]?.enable?.message}
                  </FormHelperText>
                  {(watch("addons")[index].template.isDefaultUserAddon ||
                    watch("addons")[index].enable) &&
                    field.template.requiredVars.length > 0 && (
                      <Stack spacing={2}>
                        {field.template.requiredVars?.map((required, j) => (
                          <TextField
                            key={field.id + j}
                            size="small"
                            fullWidth
                            label={required.varName}
                            defaultValue={required.defaultValue}
                            {...registerMui(
                              register(`addons.${index}.vars.${j}` as const, {
                                required: watch("addons")[index].enable,
                              })
                            )}
                            error={Boolean(errors.addons?.[index]?.vars?.[j])}
                            helperText={
                              errors.addons?.[index]?.vars?.[j] && "Required"
                            }
                          />
                        ))}
                      </Stack>
                    )}
                </React.Fragment>
              ))}
            </Stack>
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => onClose()} color="primary">
            Cancel
          </Button>
          <Button type="submit" variant="contained" color="primary">
            Confirm
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

/**
 * Context
 */
export const UserDeleteDialogContext = DialogContext<{ user: User }>(
  (props) => <UserDeleteDialog {...props} />
);
export const UserCreateDialogContext = DialogContext((props) => (
  <UserCreateDialog {...props} />
));
export const UserCreateConfirmDialogContext = DialogContext<{
  onConfirm: () => void;
  user: User;
}>((props) => <UserCreateConfirmDialog {...props} />);
