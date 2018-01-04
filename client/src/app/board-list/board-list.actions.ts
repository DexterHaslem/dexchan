import { Injectable } from "@angular/core";
import { IAppState } from "../store/root.reducer";
import { NgRedux } from "@angular-redux/store";
import { ApiService } from "../api.service";
import "rxjs/add/operator/mapTo";


@Injectable()
export class BoardListActions {
  static readonly GET_BOARDS = 'BOARD_ACTIONS_GET_BOARDS';

  constructor(private ngRedux: NgRedux<IAppState>, private api: ApiService) {

  }

  getBoards(): void {
    this.ngRedux.dispatch({
      type: BoardListActions.GET_BOARDS
    });
  }
}
