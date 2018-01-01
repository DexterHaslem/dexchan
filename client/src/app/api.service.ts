import { Injectable } from '@angular/core';
import { HttpClient } from "@angular/common/http";
import { Observable } from "rxjs/Observable";

const URLROOT = "http://localhost:8080/";

const buildUrl = (endpoint: string): string => {
  return `${URLROOT}${endpoint}`;
};

@Injectable()
export class ApiService {

  constructor(private http: HttpClient) {
  }

  getBoards(): Observable<Board[]> {
    return this.http.get<Board[]>(buildUrl('api/boards'));
  }

  getThreads(shortCode: string): Observable<Thread[]> {
    return this.http.get<Thread[]>(buildUrl(`api/${shortCode}`))
  }
}
