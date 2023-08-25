import useUrlState from "@ahooksjs/use-url-state";
import { AddTwoTone, Badge, Clear, DeleteTwoTone, ExpandLess, ExpandMore, ManageAccountsTwoTone, MoreVert, RefreshTwoTone, SearchTwoTone, Tune, ViewList, ViewModule } from "@mui/icons-material";
import { Box, Card, CardContent, CardHeader, Chip, Collapse, Divider, Fab, Grid, IconButton, InputAdornment, List, ListItem, ListItemIcon, ListItemText, Menu, MenuItem, Paper, Stack, Table, TableBody, TableCell, TableContainer, TableRow, TextField, Tooltip, Typography, useMediaQuery, useTheme } from "@mui/material";
import React, { useEffect, useState } from "react";
import { useLogin } from "../../components/LoginProvider";
import { User } from "../../proto/gen/dashboard/v1alpha1/user_pb";
import { NameAvatar } from "../atoms/NameAvatar";
import { SelectableChip } from "../atoms/SelectableChips";
import { PasswordDialogContext } from "../organisms/PasswordDialog";
import { RoleChangeDialogContext } from "../organisms/RoleChangeDialog";
import { UserCreateConfirmDialogContext, UserCreateDialogContext, UserDeleteDialogContext, UserInfoDialogContext } from "../organisms/UserActionDialog";
import { UserAddonChangeDialogContext } from "../organisms/UserAddonsChangeDialog";
import { hasAdminForRole, hasPrivilegedRole, isAdminRole, isPrivilegedRole, useUserModule } from "../organisms/UserModule";
import { UserNameChangeDialogContext } from "../organisms/UserNameChangeDialog";
import { PageTemplate } from "../templates/PageTemplate";

/**
 * view
 */
const UserMenu: React.VFC<{ user: User }> = ({ user: us }) => {
  const { loginUser } = useLogin();
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const roleChangeDialogDispatch = RoleChangeDialogContext.useDispatch();
  const userDeleteDialogDispatch = UserDeleteDialogContext.useDispatch();
  const userNameChangeDispatch = UserNameChangeDialogContext.useDispatch();
  const userAddonChangeDispatch = UserAddonChangeDialogContext.useDispatch();

  return (<>
    <Box>
      <IconButton
        color="inherit"
        disabled={loginUser?.name === us.name}
        onClick={e => setAnchorEl(e.currentTarget)}>
        <MoreVert fontSize="small" />
      </IconButton>
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={() => setAnchorEl(null)}
      >
        <MenuItem onClick={() => {
          userNameChangeDispatch(true, { user: us });
          setAnchorEl(null);
        }}>
          <ListItemIcon><Badge fontSize="small" /></ListItemIcon>
          <ListItemText>Change Name...</ListItemText>
        </MenuItem>
        <MenuItem onClick={() => {
          roleChangeDialogDispatch(true, { user: us });
          setAnchorEl(null);
        }}>
          <ListItemIcon><ManageAccountsTwoTone fontSize="small" /></ListItemIcon>
          <ListItemText>Change Role...</ListItemText>
        </MenuItem>
        <MenuItem onClick={() => {
          userAddonChangeDispatch(true, { user: us });
          setAnchorEl(null);
        }}>
          <ListItemIcon><Badge fontSize="small" /></ListItemIcon>
          <ListItemText>Change Addons...</ListItemText>
        </MenuItem>
        <MenuItem onClick={() => {
          userDeleteDialogDispatch(true, { user: us });
          setAnchorEl(null);
        }}>
          <ListItemIcon><DeleteTwoTone fontSize="small" /></ListItemIcon>
          <ListItemText>Remove User...</ListItemText>
        </MenuItem>
      </Menu>
    </Box>
  </>);
};


const UserListItem: React.VFC<{ user: User, }> = ({ user }) => {
  console.log('UserListItem');

  const [isOpen, setIsOpen] = useState(false);

  const theme = useTheme();
  const upSM = useMediaQuery(theme.breakpoints.up('sm'));

  return <Card>
    <CardHeader
      avatar={<NameAvatar name={user.displayName} />}
      title={<Stack direction='row' sx={{ mr: 2 }} onClick={() => setIsOpen(!isOpen)}>
        <Typography variant='subtitle1'>{user.name}</Typography>
        <Box sx={{ flex: '1 1 auto' }} />
        <div style={{ textAlign: 'right', maxWidth: upSM ? undefined : 100, whiteSpace: 'nowrap' }}>
          <Box component="div" sx={{ textAlign: 'right', justifyContent: "flex-end", textOverflow: 'ellipsis', overflow: 'hidden' }}>
            {user.roles && user.roles.map((v, i) => {
              return <Chip color={isPrivilegedRole(v) ? "error" : isAdminRole(v) ? "warning" : "default"} size='small' key={i} label={v} />
            })}
          </Box>
        </div>
      </Stack>}
      subheader={user.displayName}
      action={<UserMenu user={user} />}
    />
    <Collapse in={isOpen} timeout="auto" unmountOnExit>
      <CardContent>
        <Typography variant="body2">Addons</Typography>
        <List component="nav" >
          {user.addons.map((v, i) =>
            <div key={i}>
              <Divider />
              <ListItem>
                <ListItemIcon>
                  <Tune />
                </ListItemIcon>
                <ListItemText
                  disableTypography={true}
                  primary={v.template}
                  secondary={Object.keys(v.vars).length > 0 &&
                    <TableContainer component={Paper} sx={{ mt: 1 }}>
                      <Table aria-label={v.template}>
                        <TableBody>
                          {Object.keys(v.vars).map((key, j) =>
                            <TableRow key={j} sx={{ '&:last-child td, &:last-child th': { border: 0 } }} >
                              <TableCell component="th" scope="row">{key}</TableCell>
                              <TableCell align="right">{v.vars[key]}</TableCell>
                            </TableRow>
                          )}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  }
                />
              </ListItem>
            </div>
          )}
        </List>
      </CardContent>
    </Collapse>
  </Card>;
}


const UserCardItem: React.VFC<{ user: User, }> = ({ user }) => {
  console.log('UserCardItem');

  const userInfoDialogDispatch = UserInfoDialogContext.useDispatch();

  return <Card>
    <CardHeader
      avatar={<NameAvatar name={user.displayName} />}
      title={<Stack direction='row' sx={{ mr: 2 }} onClick={() => userInfoDialogDispatch(true, { user: user, defaultOpenUserAddon: true })}>
        <Typography variant='subtitle1'>{user.name}</Typography>
        <Box sx={{ flex: '1 1 auto' }} />
        <div style={{ textAlign: 'right', maxWidth: 100, whiteSpace: 'nowrap' }}>
          <Box component="div" sx={{ textAlign: 'right', justifyContent: "flex-end", textOverflow: 'ellipsis', overflow: 'hidden' }}>
            {user.roles && user.roles.map((v, i) => {
              return <Chip color={isPrivilegedRole(v) ? "error" : isAdminRole(v) ? "warning" : "default"} size='small' key={i} label={v} />
            })}
          </Box>
        </div>
      </Stack>}
      subheader={user.displayName}
      action={<UserMenu user={user} />}
    />
  </Card>;
}

const UserList: React.VFC = () => {
  const hooks = useUserModule();
  const { loginUser } = useLogin();
  const userCreateDialogDispatch = UserCreateDialogContext.useDispatch();

  const [showFilter, setShowFilter] = useState<boolean>(false);
  const [urlParam, setUrlParam] = useUrlState({
    "search": "",
    "filterRoles": [],
    "view": "list",
  }, { parseOptions: { arrayFormat: 'comma' }, stringifyOptions: { arrayFormat: 'comma', skipEmptyString: true } });

  const filterRoles: string[] = typeof urlParam.filterRoles === 'string' ? [urlParam.filterRoles] : urlParam.filterRoles;
  const pushFilterRoles = (role: string) => {
    const f = [...new Set([...filterRoles, role])].sort((a, b) => a < b ? -1 : 1)
    filterRoles && setUrlParam({ filterRoles: f });
  }
  const popFilterRoles = (role: string) => {
    const f = filterRoles.filter((v: string) => v !== role)
    filterRoles && setUrlParam({ filterRoles: f });
  }

  useEffect(() => {
    if (loginUser && !hasPrivilegedRole(loginUser.roles) && filterRoles.length === 0) {
      setUrlParam({ filterRoles: hooks.existingRoles.filter((v) => hasAdminForRole(loginUser.roles, v)) })
    }
  }, [loginUser, hooks.existingRoles.length]);

  const isUserMatchedToFilterRoles = (user: User) => {
    for (const v of user.roles) {
      for (const f of filterRoles) {
        if (v === f) {
          return true
        }
      }
    }
    return false
  }

  useEffect(() => { hooks.getUsers() }, []); // eslint-disable-line

  const isListView = urlParam.view == 'list'
  const isCardView = urlParam.view == 'card'

  return (<>
    <Paper sx={{ minWidth: 320, maxWidth: 1200, mb: 1, p: 1 }}>
      <Stack direction='row' alignItems='center' spacing={2}>
        <TextField
          InputProps={urlParam.search !== "" ? {
            startAdornment: (<InputAdornment position="start"><SearchTwoTone /></InputAdornment>),
            endAdornment: (<InputAdornment position="end">
              <IconButton size="small" tabIndex={-1} onClick={() => { setUrlParam({ search: "" }) }} >
                <Clear />
              </IconButton>
            </InputAdornment>)
          } : {
            startAdornment: (<InputAdornment position="start"><SearchTwoTone /></InputAdornment>),
          }}
          placeholder="Search"
          size='small'
          value={urlParam.search}
          onChange={e => setUrlParam({ search: e.target.value })}
          sx={{ flexGrow: 0.5 }}
        />
        <Box sx={{ flexGrow: 1 }} />
        {isListView && <Tooltip title="CardView" placement="top">
          <IconButton color="inherit" onClick={() => { setUrlParam({ view: "card" }) }}>
            <ViewModule />
          </IconButton>
        </Tooltip>}
        {isCardView && <Tooltip title="ListView" placement="top">
          <IconButton color="inherit" onClick={() => { setUrlParam({ view: "list" }) }}>
            <ViewList />
          </IconButton>
        </Tooltip>}
        <Tooltip title="Refresh" placement="top">
          <IconButton color="inherit" onClick={() => { hooks.getUsers() }}>
            <RefreshTwoTone />
          </IconButton>
        </Tooltip>
        <Tooltip title="Add new user" placement="top">
          <Fab size='small' color='primary' onClick={() => userCreateDialogDispatch(true)} sx={{ flexShrink: 0 }} >
            <AddTwoTone />
          </Fab>
        </Tooltip>
      </Stack >
      <Box component="div" sx={{ justifyContent: "flex-end", textOverflow: 'ellipsis', overflow: 'hidden' }} >
        <IconButton size="small" color="inherit" onClick={() => { setShowFilter(!showFilter) }}>
          {showFilter ? < ExpandLess /> : <ExpandMore />}
        </IconButton>
        <Typography color="text.secondary" variant="caption">Filter by Roles</Typography>
        {filterRoles.length > 0 &&
          <Grid container sx={{ pt: 1 }}>
            {filterRoles.map((v, i) =>
              <SelectableChip key={v} label={v} sx={{ m: 0.1 }}
                color={isPrivilegedRole(v) ? "error" : isAdminRole(v) ? "warning" : "default"}
                defaultChecked={true} onChecked={() => { popFilterRoles(v) }} />
            )}
          </Grid>}
        <Collapse in={showFilter} timeout="auto" unmountOnExit sx={{ pt: 1 }}>
          <Divider />
          <Typography color="text.secondary" variant="caption">Existing Roles</Typography>
          <Grid container sx={{ pt: 1 }}>
            {hooks.existingRoles.map((v, i) =>
              <SelectableChip key={v} label={v} sx={{ m: 0.1 }}
                color={isPrivilegedRole(v) ? "error" : isAdminRole(v) ? "warning" : "default"}
                checked={filterRoles?.includes(v)} onChecked={(checked) => { checked ? pushFilterRoles(v) : popFilterRoles(v) }} />
            )}
          </Grid>
        </Collapse>
      </Box>

    </Paper>
    {
      !hooks.users.filter((us) => urlParam.search === '' || Boolean(us.name.match(urlParam.search))).length &&
      <Paper sx={{ minWidth: 320, maxWidth: 1200, mb: 1, p: 4 }}>
        <Typography variant='subtitle1' sx={{ color: 'text.secondary', textAlign: 'center' }}>No Users found.</Typography>
      </Paper>
    }
    <Grid container spacing={0.5}>
      {hooks.users
        .filter((us) => urlParam.search === '' || Boolean(us.name.match(urlParam.search)))
        .filter((us) => us.status === 'Active')
        .filter((us) => (filterRoles.length == 0 || isUserMatchedToFilterRoles(us)))
        .map((us) =>
          <Grid item key={us.name} xs={12} sm={isListView ? 12 : 6} md={isListView ? 12 : 6} lg={isListView ? 12 : 4}>
            {isListView && <UserListItem user={us} />}
            {isCardView && <UserCardItem user={us} />}
          </Grid>
        )}
    </Grid >
  </>);
};

export const UserPage: React.VFC = () => {
  console.log('UserPage');

  return (
    <PageTemplate title="Users">
      <PasswordDialogContext.Provider>
        <UserCreateConfirmDialogContext.Provider>
          <UserCreateDialogContext.Provider>
            <RoleChangeDialogContext.Provider>
              <UserAddonChangeDialogContext.Provider>
                <UserDeleteDialogContext.Provider>
                  <UserInfoDialogContext.Provider>
                    <UserList />
                  </UserInfoDialogContext.Provider>
                </UserDeleteDialogContext.Provider>
              </UserAddonChangeDialogContext.Provider>
            </RoleChangeDialogContext.Provider>
          </UserCreateDialogContext.Provider>
        </UserCreateConfirmDialogContext.Provider>
      </PasswordDialogContext.Provider>
    </PageTemplate>
  );
}
