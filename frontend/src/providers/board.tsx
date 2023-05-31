'use client';

import React, { createContext, useContext, useMemo, useReducer } from 'react';
import { BoardWithMembers } from '../api/boards';

type State = {
  board: BoardWithMembers | null;
  isOwner: boolean;
};

type Action = {
  type: 'set_board' | 'set_is_owner';
  payload?: any;
};

type Dispatch = (action: Action) => void;

export type BoardContextType = {
  state: State;
  dispatch: Dispatch;
};

const INITIAL_STATE = {
  board: null,
  isOwner: false,
};

const boardReducer = (state: State, action: Action) => {
  switch (action.type) {
    case 'set_board':
      return { ...state, board: action.payload };
    case 'set_is_owner':
      return { ...state, isOwner: action.payload };
    default: {
      throw new Error(`Unhandled action type: ${action.type}`);
    }
  }
};

const BoardContext = createContext<BoardContextType>({ state: INITIAL_STATE, dispatch: () => null });

export const BoardProvider = ({ children }: { children: React.ReactNode }) => {
  const [state, dispatch] = useReducer(boardReducer, INITIAL_STATE);

  const value = useMemo(() => ({ state, dispatch }), [state]);
  return <BoardContext.Provider value={value}>{children}</BoardContext.Provider>;
};

export const useBoard = () => {
  return useContext(BoardContext);
};
