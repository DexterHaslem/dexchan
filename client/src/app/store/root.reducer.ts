import { Reducer } from 'redux';
import { appReducer } from "./app.reducer";

export interface IAction {
  type: string;
  payload?: any;
}

export interface IAppState {
  boards: Board[];
  selectedBoard: string;

  threads: Thread[];
  selectedThreadID: number;

  posts: Post[];
}

export const INITIAL_STATE: IAppState = {
  boards: [],
  posts: [],
  selectedBoard: '',
  selectedThreadID: 0,
  threads: []
};

// export const rootReducer: Reducer<IAppState> = combineReducers<IAppState>({
//
// });

export const rootReducer: Reducer<IAppState> = appReducer;
