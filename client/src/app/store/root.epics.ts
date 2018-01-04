import { Injectable } from '@angular/core';
import { combineEpics } from 'redux-observable';

@Injectable()
export class RootEpics {
  constructor() {
  }

  public combineEpics() {
    return combineEpics(

    );
  }
}
