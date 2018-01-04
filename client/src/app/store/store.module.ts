import { NgModule } from '@angular/core';
import { DevToolsExtension, NgRedux, NgReduxModule } from '@angular-redux/store';
import { createLogger } from 'redux-logger';
import { createEpicMiddleware } from 'redux-observable';
import { IAppState, INITIAL_STATE, rootReducer } from './root.reducer';
import { environment } from '../../environments/environment';
import { RootEpics } from "./root.epics";

@NgModule({
  imports: [NgReduxModule],
  providers: [RootEpics]
})
export class StoreModule {
  constructor(public store: NgRedux<IAppState>,
              devTools: DevToolsExtension,
              rootEpics: RootEpics) {

    let middleware = [createEpicMiddleware(rootEpics.combineEpics())];
    if (!environment.production) {
      middleware = [createLogger(), ...middleware];
    }
    store.configureStore(
      rootReducer,
      INITIAL_STATE,
      middleware,
      [devTools.isEnabled() ? devTools.enhancer() : f => f]
    );
  }
}
