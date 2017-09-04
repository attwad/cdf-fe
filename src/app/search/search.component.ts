import { Component, OnInit } from '@angular/core';
import {Router} from '@angular/router';
import {PageEvent} from '@angular/material';
import {LessonsService, SearchResponse} from '../lessons.service';

import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';

// Observable class extensions
import 'rxjs/add/observable/of';

// Observable operators
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/debounceTime';
import 'rxjs/add/operator/distinctUntilChanged';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.css']
})
export class SearchComponent implements OnInit {
  searchResponse?: SearchResponse;
  private searchQueries = new Subject<string>();
  loading = false;
  error?: string;
  private query?: string;
  // MdPaginator Inputs
  total = 0;
  pageSize = 5;
  pageSizeOptions = [5, 10, 25, 50];

  constructor(
      private lessonsService: LessonsService,
      private router: Router) { }

  // Push a search term into the observable stream.
  search(query: string): void {
    if (!query || query.length < 2) {
      this.searchResponse = null;
      this.error = null;
      this.loading = false;
      return;
    }
    this.searchQueries.next(query);
  }

  onError(error: any): Observable<SearchResponse> {
    console.error(error);
    this.error = error.message;
    return Observable.of<SearchResponse>(null);
  }

  onSuccess(response: SearchResponse): void {
    if (response) {
      this.error = null;
    }
    this.loading = false;
    this.total = response.total;
    this.searchResponse = response;
  }

  onPageEvent(pageEvent: PageEvent) {
    this.pageSize = pageEvent.pageSize;
    this.lessonsService.search(
      this.query,
      pageEvent.pageIndex * pageEvent.pageSize, pageEvent.pageSize)
      .catch(error => this.onError(error))
      .subscribe((response: SearchResponse) => this.onSuccess(response));
  }

  ngOnInit() {
    this.searchQueries
    // wait 300ms after each keystroke before considering the term
    .debounceTime(300)
    // ignore if next search term is same as previous
    .distinctUntilChanged()
    // switch to new observable each time the term changes
    .switchMap(query => {
      this.loading = true;
      this.query = query;
      // TODO: Why no more queries are issued after an error occurs?
      return this.lessonsService.search(query, 0, this.pageSize);
    })
    .catch(error => this.onError(error))
    .subscribe((response: SearchResponse) => this.onSuccess(response));
  }

}
