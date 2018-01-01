import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap, Router } from "@angular/router";
import { of } from "rxjs/observable/of";

@Component({
  selector: 'app-posts',
  templateUrl: './posts.component.html',
  styleUrls: ['./posts.component.css']
})
export class PostsComponent implements OnInit {
  params: ParamMap;

  constructor(private route: ActivatedRoute,
              private router: Router) {

  }

  ngOnInit() {
    // TODO: api call. had to force subscribe to get the cold observable to fire
    this.route.paramMap.switchMap((params: ParamMap) => {
      return of(params);
    }).subscribe(v => this.params = v);
  }

}
