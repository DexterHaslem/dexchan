import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import 'rxjs/add/operator/switchMap';
import { of } from "rxjs/observable/of";

@Component({
  selector: 'app-threads',
  templateUrl: './threads.component.html',
  styleUrls: ['./threads.component.css']
})
export class ThreadsComponent implements OnInit {

  public boardName: string;

  constructor(private route: ActivatedRoute,
              private router: Router) {

  }

  ngOnInit() {
    // TODO: api call. had to force subscribe to get the cold observable to fire
    this.route.paramMap.switchMap((params: ParamMap) => {
      return of(params.get('board'));
    }).subscribe(v => this.boardName = v);
  }

}
