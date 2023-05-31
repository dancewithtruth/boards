'use client';

import React, { createContext, useContext, useEffect, useMemo, useReducer } from 'react';
import { User, getUserByJwt } from '@/api/users';
import { LOCAL_STORAGE_AUTH_TOKEN } from '@/constants';

type State = {
  user: User | null;
  isAuthenticated: boolean;
};

type Action = {
  type: 'set_user' | 'set_is_authenticated';
  payload?: any;
};

type Dispatch = (action: Action) => void;

export type UserContextType = {
  state: State;
  dispatch: Dispatch;
};

const INITIAL_STATE = {
  user: null,
  isAuthenticated: false,
};

const userReducer = (state: State, action: Action) => {
  switch (action.type) {
    case 'set_user':
      return { ...state, user: action.payload };
    case 'set_is_authenticated':
      return { ...state, isAuthenticated: action.payload };
    default: {
      throw new Error(`Unhandled action type: ${action.type}`);
    }
  }
};

const UserContext = createContext<UserContextType>({ state: INITIAL_STATE, dispatch: () => null });

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const [state, dispatch] = useReducer(userReducer, INITIAL_STATE);

  const authenticate = async () => {
    const authToken = localStorage.getItem(LOCAL_STORAGE_AUTH_TOKEN);
    if (authToken) {
      try {
        const user = await getUserByJwt(authToken);
        dispatch({ type: 'set_user', payload: user });
        dispatch({ type: 'set_is_authenticated', payload: true });
      } catch (e) {
        localStorage.removeItem(LOCAL_STORAGE_AUTH_TOKEN);
      }
    }
  };

  useEffect(() => {
    authenticate();
  }, []);

  const value = useMemo(() => ({ state, dispatch }), [state]);
  return <UserContext.Provider value={value}>{children}</UserContext.Provider>;
};

export const useUser = () => {
  return useContext(UserContext);
};
