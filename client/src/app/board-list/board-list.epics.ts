import { Injectable } from "@angular/core";
import { ApiService } from "../api.service";
import { ActionsObservable } from "redux-observable";
import { BoardListActions } from "./board-list.actions";
import { IAction } from "../store/root.reducer";

@Injectable()
export class BoardListEpics {

  getBoards = (action$: ActionsObservable<IAction>) => {
    return action$.ofType(BoardListActions.GET_BOARDS)
      .mergeMap(({payload}) => {
        return this.api.getBoards();
      });
  }

  constructor(private api: ApiService) {
  }
}
