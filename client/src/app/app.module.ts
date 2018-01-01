import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';


import { AppComponent } from './app.component';
import { BoardListComponent } from './board-list/board-list.component';
import { ThreadsComponent } from './threads/threads.component';
import { PostsComponent } from './posts/posts.component';
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { ApiService } from './api.service';
import { AppRoutingModule } from './app-routing.module';
import {
  MatButtonModule, MatCardModule, MatCommonModule, MatGridListModule, MatIconModule, MatListModule,
  MatOptionModule, MatPaginatorModule, MatSelectModule, MatSlideToggleModule, MatSnackBarModule
} from "@angular/material";
import { HttpClientModule } from "@angular/common/http";
import { BoardInfoCardComponent } from './board-info-card/board-info-card.component';
import { ThreadItemComponent } from './thread-item/thread-item.component';


@NgModule({
  declarations: [
    AppComponent,
    BoardListComponent,
    ThreadsComponent,
    PostsComponent,
    BoardInfoCardComponent,
    ThreadItemComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    AppRoutingModule,
    MatCommonModule,
    MatSlideToggleModule,
    MatGridListModule,
    MatSnackBarModule,
    MatCardModule,
    MatListModule,
    MatButtonModule,
    MatIconModule,
    MatPaginatorModule,
    MatSelectModule,
    MatOptionModule,
    HttpClientModule
  ],
  providers: [ApiService],
  bootstrap: [AppComponent]
})
export class AppModule {
}
