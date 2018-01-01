import { Component, OnInit } from '@angular/core';
import { ApiService } from "../api.service";
import { Observable } from "rxjs/Observable";

@Component({
  selector: 'app-board-list',
  templateUrl: './board-list.component.html',
  styleUrls: ['./board-list.component.css']
})
export class BoardListComponent implements OnInit {

  boards$: Observable<Board[]>;

  constructor(private api: ApiService) {
  }

  ngOnInit() {
    this.boards$ = this.api.getBoards();
  }

}
