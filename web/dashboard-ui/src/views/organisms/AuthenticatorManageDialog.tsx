import { Delete } from "@mui/icons-material";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Grid,
  IconButton,
  Stack,
  Typography
} from "@mui/material";
import Box from '@mui/material/Box';
import { useSnackbar } from 'notistack';
import { useEffect, useState } from "react";
import { base64url } from '../../components/Base64';
import { DialogContext } from "../../components/ContextProvider";
import { User } from "../../proto/gen/dashboard/v1alpha1/user_pb";
import { Credential } from "../../proto/gen/dashboard/v1alpha1/webauthn_pb";
import { useWebAuthnService } from '../../services/DashboardServices';
/**
 * view
 */

export const AuthenticatorManageDialog: React.VFC<{ onClose: () => void, user: User }> = ({ onClose, user }) => {
  console.log('AuthenticatorManageDialog');
  const webauthnService = useWebAuthnService();
  const { enqueueSnackbar } = useSnackbar();

  const [credentials, setCredentials] = useState<Credential[]>([]);

  const registerdCredId = localStorage.getItem(`credId`)
  const isRegistered = Boolean(registerdCredId && credentials.map(c => c.id).includes(registerdCredId!));

  const [isWebAuthnAvailable, setIsWebAuthnAvailable] = useState(false);

  const checkWebAuthnAvailable = () => {
    if (window.PublicKeyCredential) {
      PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable()
        .then(uvpaa => { setIsWebAuthnAvailable(uvpaa) });
    }
  }
  useEffect(() => { checkWebAuthnAvailable() }, []);

  console.log("credId", registerdCredId, "isRegistered", isRegistered, "isWebAuthnAvailable", isWebAuthnAvailable);

  /**
   * listCredentials
  */
  const listCredentials = async () => {
    console.log("listCredentials");
    try {
      const resp = await webauthnService.listCredentials({ userName: user.name });
      setCredentials(resp.credentials);
      console.log(resp.credentials);
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

      // Credential is allowed to access only id and type so use any.
      const cred: any = await navigator.credentials.create(opt);
      if (cred === null) {
        console.log("cred is null");
        throw Error('credential is null');
      }

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
      localStorage.setItem(`credId`, credential.rawId);

      const finResp = await webauthnService.finishRegistration({ userName: user.name, credentialCreationResponse: JSON.stringify(credential) });
      enqueueSnackbar(finResp.message, { variant: 'success' });
      listCredentials();
    }
    catch (error) {
      handleError(error);
    }
  }

  /**
   * removeCredentials
  */
  const removeCredentials = async (id: string) => {
    console.log("removeCredentials");
    if (!confirm("remove?")) { return }
    try {
      const resp = await webauthnService.deleteCredential({ userName: user.name, credId: id });
      enqueueSnackbar(resp.message, { variant: 'success' });
      listCredentials();
      if (id === registerdCredId) {
        localStorage.removeItem(`credId`);
      }
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
    error instanceof DOMException || msg && enqueueSnackbar(msg, { variant: 'error' });
  }

  return (
    <Dialog open={true}
      fullWidth maxWidth={'sm'}>
      <DialogTitle>WebAuthn Credentials</DialogTitle>
      <DialogContent>
        <Box sx={{ p: 2, border: '1px grey', borderRadius: '4px' }}>
          <Stack alignItems="center" >
            {credentials.length === 0 && <Typography>No credentials</Typography>}
            {credentials.length > 0 && <Grid container sx={{ p: 2 }}>
              {credentials.map((field, index) => {
                return (
                  <>
                    <Grid item xs={11} sx={{ m: 'auto' }} key={index}>
                      <Typography sx={registerdCredId === field.id && {
                        color: 'red',
                      } || undefined}>{field.id}</Typography>
                      <Typography variant="caption" display="block">{field.displayName}</Typography>
                    </Grid>
                    <Grid item xs={1} sx={{ m: 'auto', textAlign: 'end' }}>
                      <IconButton edge="end" aria-label="delete" onClick={() => { removeCredentials(field.id) }}><Delete /></IconButton>
                    </Grid>
                  </>
                )
              })}
            </Grid>}
          </Stack>
          {/* <List>
            {credentials.map((field, index) =>
              <div key={index}>
                <ListItem
                  secondaryAction={<IconButton edge="end" aria-label="delete" onClick={() => { removeCredentials(field.id) }}><Delete /></IconButton>}
                >
                  <ListItemText
                    primary={field.id}
                    primaryTypographyProps={registerdCredId === field.id && {
                      color: 'red',
                    } || undefined}

                    secondary={<Typography variant="caption" display="block">{field.displayName}</Typography>}
                  />
                </ListItem>
                <Divider />
              </div>
            )}
          </List> */}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={() => onClose()} color="primary">Close</Button>
        {!isRegistered && isWebAuthnAvailable
          ? <Button onClick={() => registerNewAuthenticator()} variant="contained" color="secondary">Register</Button>
          : undefined}
      </DialogActions>
    </Dialog >
  );
};

/**
 * Context
 */
export const AuthenticatorManageDialogContext = DialogContext<{ user: User }>(
  props => (<AuthenticatorManageDialog {...props} />));
