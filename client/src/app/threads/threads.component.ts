import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';
import 'rxjs/add/operator/switchMap';
import { ApiService } from "../api.service";
import { Observable } from "rxjs/Observable";

@Component({
  selector: 'app-threads',
  templateUrl: './threads.component.html',
  styleUrls: ['./threads.component.css']
})
export class ThreadsComponent implements OnInit {

  threads$: Observable<Thread[]>;

  constructor(private route: ActivatedRoute,
              private api: ApiService) {
  }

  foods = [
    {value: 'steak-0', viewValue: 'Steak'},
    {value: 'pizza-1', viewValue: 'Pizza'},
    {value: 'tacos-2', viewValue: 'Tacos'}
  ];

  ngOnInit() {
    this.threads$ = this.route.paramMap
      .switchMap((params: ParamMap) => this.api.getThreads(params.get('board')));
  }
}
