import { Component, OnInit } from '@angular/core';
import {Router} from '@angular/router';
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

  ngOnInit() {
    this.searchQueries
    // wait 300ms after each keystroke before considering the term
    .debounceTime(300)
    // ignore if next search term is same as previous
    .distinctUntilChanged()
    // switch to new observable each time the term changes
    .switchMap(query => {
      this.loading = true;
      // TODO: Why no more queries are issued after an error occurs?
      return this.lessonsService.search(query);
    })
    .catch(error => {
      console.error(error);
      this.error = error.message;
      return Observable.of<SearchResponse>(null);
    })
    .subscribe((response: SearchResponse) => {
      if (response) {
        this.error = null;
      }
      this.loading = false;
      this.searchResponse = response;
    });
  }

}
