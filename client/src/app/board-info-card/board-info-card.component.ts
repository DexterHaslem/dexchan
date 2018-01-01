import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-board-info-card',
  templateUrl: './board-info-card.component.html',
  styleUrls: ['./board-info-card.component.css']
})
export class BoardInfoCardComponent implements OnInit {

  @Input()
  board: Board;

  constructor() {

  }

  getRouterLink(): string {
    return `${this.board.shortCode}`;
  }

  ngOnInit() {
  }

}
