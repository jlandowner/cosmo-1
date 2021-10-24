import {
  Box, Chip, colors, Container, CssBaseline, Divider, IconButton,
  Link, ListItemIcon, ListItemText,
  Menu, MenuItem, Stack, Toolbar, Typography
} from "@mui/material";
import MuiAppBar, { AppBarProps } from '@mui/material/AppBar';
import { experimentalStyled as styled } from "@mui/material/styles";
import { AccountCircle, ExitToApp, LockOutlined, Menu as MenuIcon, SupervisorAccountTwoTone, VpnKey, WebTwoTone } from "@mui/icons-material";
import React from "react";
import { ErrorBoundary } from 'react-error-boundary';
import { Link as RouterLink } from "react-router-dom";
import { useLogin } from "../../components/LoginProvider";
import logo from "../../logo-with-name-small.png";
import { NameAvatar } from "../atoms/NameAvatar";
import { PasswordChangeDialogContext } from "../organisms/PasswordChangeDialog";


const AppBar = styled(MuiAppBar)<AppBarProps>(({ theme }) => ({
  zIndex: theme.zIndex.drawer + 1,
  transition: theme.transitions.create(['width', 'margin'], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
}));

const Copyright = () => {
  return (
    <Typography variant="body2" color="textSecondary" align="center">
      {"Copyright © "}
      <Link href="https://github.com/cosmo-workspace">
        cosmo-workspace
      </Link>{` ${new Date().getFullYear()}.`}
    </Typography>
  );
};

interface PageTemplateProps {
  children: React.ReactNode;
  title: string;
}

export const PageTemplate: React.FC<PageTemplateProps> = ({ children, title, }) => {

  const { loginUser, logout } = useLogin();
  const passwordChangeDialogDispach = PasswordChangeDialogContext.useDispatch();
  const isAdmin = (loginUser?.role === 'cosmo-admin');
  const isSignIn = Boolean(loginUser);

  const changePassword = () => {
    console.log('changePassword');
    passwordChangeDialogDispach(true);
    setAnchorEl(null);
  }

  const [menuAnchorEl, setMenuAnchorEl] = React.useState<null | HTMLElement>(null);
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);


  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <AppBar position="absolute">
        <Toolbar sx={{ pr: 3 }} >
          <IconButton
            edge="start"
            color="inherit"
            onClick={(e) => setMenuAnchorEl(e.currentTarget)}
            sx={{ mr: 2 }}>
            <MenuIcon />
          </IconButton>
          <Menu
            id="basic-menu"
            anchorEl={menuAnchorEl}
            open={Boolean(menuAnchorEl)}
            onClose={() => setMenuAnchorEl(null)}
          >
            {loginUser && <RouterLink to="/workspace" style={{ textDecoration: "none", color: "inherit" }}>
              <MenuItem>
                <ListItemIcon>
                  <WebTwoTone />
                </ListItemIcon>
                <ListItemText primary="Workspaces" />
              </MenuItem>
            </RouterLink>}
            {loginUser && isAdmin && <RouterLink to="/user" style={{ textDecoration: "none", color: "inherit" }}>
              <MenuItem>
                <ListItemIcon>
                  <SupervisorAccountTwoTone />
                </ListItemIcon>
                <ListItemText primary="Users" />
              </MenuItem>
            </RouterLink>}
            {!loginUser && <RouterLink to="/signin" style={{ textDecoration: "none", color: "inherit" }}>
              <MenuItem>
                <ListItemIcon>
                  <LockOutlined />
                </ListItemIcon>
                <ListItemText primary="sign in" />
              </MenuItem>
            </RouterLink>}
          </Menu>
          <RouterLink to="/workspace" style={{ textDecoration: "none", color: "inherit" }}>
            <img alt="cosmo" src={logo} height={40} />
          </RouterLink>
          <Box sx={{ flexGrow: 1 }} />
          <Box>
            <IconButton
              color="inherit"
              onClick={(e) => setAnchorEl(e.currentTarget)}
              disabled={!isSignIn}>
              <AccountCircle fontSize="large" />
            </IconButton>
            <Menu
              id="basic-menu"
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={() => setAnchorEl(null)}
            >
              <Stack alignItems="center" spacing={1} sx={{ mt: 1, mb: 2 }}>
                <NameAvatar name={loginUser?.displayName} sx={{ width: 50, height: 50 }} />
                <Typography>{loginUser?.displayName}</Typography>
                <Typography color={colors.grey[700]} fontSize="small">{loginUser?.id}</Typography>
                {loginUser?.role && <Chip variant="outlined" size="small" label={loginUser?.role} />}
              </Stack>
              <Divider sx={{ mb: 1 }} />
              {isSignIn && <MenuItem onClick={() => changePassword()}>
                <ListItemIcon><VpnKey fontSize="small" /></ListItemIcon>
                <ListItemText>Change Password...</ListItemText>
              </MenuItem>}
              <Divider />
              {isSignIn && <MenuItem onClick={() => logout()}>
                <ListItemIcon><ExitToApp fontSize="small" /></ListItemIcon>
                <ListItemText>Sign out</ListItemText>
              </MenuItem>}
            </Menu>
          </Box>

        </Toolbar>
      </AppBar>

      <Box
        component="main"
        sx={{
          backgroundColor: (theme) =>
            theme.palette.mode === 'light'
              ? theme.palette.grey[100]
              : theme.palette.grey[900],
          flexGrow: 1,
          height: '100vh',
          overflow: 'auto',
        }}
      >
        <Toolbar />
        <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
          <Typography
            component="h2"
            variant="h5"
            color="inherit"
            noWrap
            sx={{ mb: 1 }}
          >
            {title}
          </Typography>
          <ErrorBoundary
            FallbackComponent={
              ({ error, resetErrorBoundary }) => {
                return (
                  <div>
                    <p>Something went wrong:</p>
                    <pre>{error.message}</pre>
                  </div>
                )
              }
            }
          >
            {children}
          </ErrorBoundary>
          <Box pt={4}>
            <Copyright />
          </Box>
        </Container>
      </Box>
    </Box>
  );
};