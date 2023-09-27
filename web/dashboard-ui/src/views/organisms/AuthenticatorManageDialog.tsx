import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  Typography
} from "@mui/material";
import { useSnackbar } from 'notistack';
import { useEffect, useState } from "react";
import { base64url } from '../../components/Base64';
import { DialogContext } from "../../components/ContextProvider";
import { User } from "../../proto/gen/dashboard/v1alpha1/user_pb";
import { useWebAuthnService } from '../../services/DashboardServices';

/**
 * view
 */

export const AuthenticatorManageDialog: React.VFC<{ onClose: () => void, user: User }> = ({ onClose, user }) => {
  console.log('AuthenticatorManageDialog');
  const webauthnService = useWebAuthnService();
  const { enqueueSnackbar } = useSnackbar();

  const [credentials, setCredentials] = useState<string[]>([]);

  /**
   * listCredentials
  */
  const listCredentials = async () => {
    console.log("listCredentials");
    try {
      const resp = await webauthnService.listCredentials({ userName: user.name });
      setCredentials(resp.credentials.map(cred => cred.id));
    }
    catch (error) {
      handleError(error);
    }
  }
  useEffect(() => { listCredentials() }, []);

  /**
   * registerNewAuthenticator
   */
  const registerNewAuthenticator = async () => {
    try {
      const resp = await webauthnService.beginRegistration({ userName: user.name });
      const options = JSON.parse(resp.credentialCreationOptions);

      const opt: CredentialCreationOptions = JSON.parse(JSON.stringify(options));
      if (options.publicKey?.user.id) {
        opt.publicKey!.user.id = base64url.decode(options.publicKey?.user.id);
      }
      if (options.publicKey?.challenge) {
        opt.publicKey!.challenge = base64url.decode(options.publicKey?.challenge);
      }
      console.log("opt", opt);

      const cred = await navigator.credentials.create(opt);
      if (cred === null) {
        console.log("cred is null");
        throw Error('credential is null');
      }
      localStorage.setItem(`credId`, cred.id);

      const credential = {
        id: cred.id,
        rawId: base64url.encode(cred.rawId),
        type: cred.type,
        response: {
          clientDataJSON: base64url.encode(cred.response.clientDataJSON),
          attestationObject: base64url.encode(cred.response.attestationObject)
        }
      };
      console.log("credential", credential);

      const finResp = await webauthnService.finishRegistration({ userName: user.name, credentialCreationResponse: JSON.stringify(credential) });
      enqueueSnackbar(finResp.message, { variant: 'success' });
      listCredentials();
    }
    catch (error) {
      handleError(error);
    }
  }

  /**
   * error handler
   */
  const handleError = (error: any) => {
    console.log(error);
    const msg = error?.message;
    msg && enqueueSnackbar(msg, { variant: 'error' });
  }

  return (
    <Dialog open={true}
      fullWidth maxWidth={'sm'}>
      <DialogTitle>Authenticators</DialogTitle>
      <DialogContent>
        <Stack spacing={3}>
          {credentials.map((field, index) =>
            <Typography key={index}>{field}</Typography>
          )}
        </Stack>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => onClose()} color="primary">Close</Button>
        <Button onClick={() => registerNewAuthenticator()} variant="contained" color="secondary">Register</Button>
      </DialogActions>
    </Dialog>
  );
};

/**
 * Context
 */
export const AuthenticatorManageDialogContext = DialogContext<{ user: User }>(
  props => (<AuthenticatorManageDialog {...props} />));
