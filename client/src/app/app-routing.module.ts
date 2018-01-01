import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from "@angular/router";
import { ThreadsComponent } from "./threads/threads.component";
import { PostsComponent } from "./posts/posts.component";
import { BoardListComponent } from "./board-list/board-list.component";

const routes: Routes = [
  {path: ':board', component: ThreadsComponent, pathMatch: 'full'},
  {path: ':board/:thread', component: PostsComponent, pathMatch: 'full'},
  {path: '', component: BoardListComponent, pathMatch: 'full'}
];

@NgModule({
  imports: [
    CommonModule,
    RouterModule.forRoot(routes)
  ],
  exports: [
    RouterModule
  ],
  declarations: []
})
export class AppRoutingModule {
}
