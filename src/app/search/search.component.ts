import { Component, OnInit } from '@angular/core';
import {Router} from '@angular/router';
import {LessonsService, SearchResponse} from '../lessons.service';

import { Observable }        from 'rxjs/Observable';
import { Subject }           from 'rxjs/Subject';

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

  constructor(
      private lessonsService: LessonsService,
      private router: Router) { }

  // Push a search term into the observable stream.
  search(query: string): void {
    this.searchQueries.next(query);
    console.log(this.searchQueries);
  }

  ngOnInit() {
    this.searchQueries
    .debounceTime(300)         // wait 300ms after each keystroke before considering the term
    .distinctUntilChanged()    // ignore if next search term is same as previous
    .switchMap(query => query  // switch to new observable each time the term changes
      ? this.lessonsService.search(query)
      : Observable.of<SearchResponse>(null))
    .catch(error => {
      // TODO: add real error handling
      console.error(error);
      return Observable.of<SearchResponse>(null);
    })
    .subscribe((response: SearchResponse) => this.searchResponse = response);
  }

}
