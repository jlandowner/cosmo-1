import useUrlState from "@ahooksjs/use-url-state";
import { Timestamp } from "@bufbuild/protobuf";
import { useSnackbar } from "notistack";
import { useEffect, useState } from "react";
import { ModuleContext } from "../../components/ContextProvider";
import { useHandleError, useLogin } from "../../components/LoginProvider";
import { Event } from "../../proto/gen/dashboard/v1alpha1/event_pb";
import { User } from "../../proto/gen/dashboard/v1alpha1/user_pb";
import { useUserService } from "../../services/DashboardServices";
import {
  isAdminUser,
  setUsersFuncFilteredByAccesibleRoles,
} from "./UserModule";
/**
 * hooks
 */
const useEvent = () => {
  console.log("useEvent");
  const { enqueueSnackbar } = useSnackbar();
  const { handleError } = useHandleError();

  const [events, setEvents] = useState<Event[]>([]);
  const userService = useUserService();

  const { loginUser, updateClock } = useLogin();
  const isAdmin = isAdminUser(loginUser);
  const [users, setUsers] = useState<User[]>([loginUser || new User()]);

  const [urlParam, setUrlParam] = useUrlState(
    { search: "", user: loginUser?.name || "" },
    {
      stringifyOptions: { skipEmptyString: true },
    }
  );
  const search: string = urlParam.search || "";
  const setSearch = (word: string) => setUrlParam({ search: word });

  const userName: string = urlParam.user || loginUser?.name;
  const user = users.find((u) => u.name === userName) || new User();
  const setUser = (name: string) => setUrlParam({ user: name });

  useEffect(() => {
    getEvents(userName);
    isAdmin && getUsers();
  }, [userName]); // eslint-disable-line

  useEffect(() => {
    if (users.length > 1) {
      users.find((u) => u.name === userName) ||
        enqueueSnackbar(`User ${userName} is not found`, {
          variant: "error",
        });
    }
  }, [users]); // eslint-disable-line

  const getUsers = async () => {
    console.log("useEvent:getUsers");
    try {
      const result = await userService.getUsers({});
      setUsers(setUsersFuncFilteredByAccesibleRoles(result.items, loginUser));
    } catch (error) {
      handleError(error);
    }
  };

  const getEvents = async (userName: string) => {
    console.log("useEvent:getEvents");
    try {
      const result = await userService.getEvents({ userName: userName });
      setEvents(result.items);
      updateClock();
    } catch (error) {
      handleError(error);
    }
  };

  return {
    search,
    setSearch,
    user,
    setUser,
    users,
    setUsers,
    events,
    setEvents,
    getEvents,
    getUsers,
    userName,
  };
};

export function getTime(timestamp?: Timestamp): number {
  if (!timestamp) {
    return 0;
  }
  return timestamp.toDate().getTime();
}

export function latestTime(event?: Event): number {
  if (!event) {
    return 0;
  }
  const t1 = getTime(event.series?.lastObservedTime);
  const t2 = getTime(event.eventTime);
  return t1 > t2 ? t1 : t2;
}

export const EventContext = ModuleContext(useEvent);
export const useEventModule = EventContext.useContext;
