import { Component, OnInit, Input} from '@angular/core';
import { Title } from '@angular/platform-browser';
import {Router, ActivatedRoute, ParamMap} from '@angular/router';
import {PageEvent} from '@angular/material';
import {LessonsService, SearchResponse} from '../lessons.service';

import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';

// Observable class extensions
import 'rxjs/add/observable/of';

// Observable operators
import 'rxjs/add/operator/catch';

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
  query?: string;
  // MdPaginator Inputs
  total = 0;
  pageSize = 5;
  pageSizeOptions = [5, 10, 25, 50];

  constructor(
      title: Title,
      private lessonsService: LessonsService,
      private router: Router,
      private route: ActivatedRoute) {
        title.setTitle('Search College de France audio transcripts');
      }

  search(query: string): void {
    this.router.navigate(['/search', {q: query}]);
  }

  onError(error: any): Observable<SearchResponse> {
    console.error(error);
    this.error = error.message;
    return Observable.of<SearchResponse>(null);
  }

  onSuccess(response: SearchResponse): void {
    if (response) {
      this.error = null;
      this.total = response.total;
    }
    this.loading = false;
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
    this.route.paramMap
      .switchMap((params: ParamMap) => {
        this.query = params.get('q');
        if (this.query) {
          this.loading = true;
          return this.lessonsService.search(this.query, 0, this.pageSize);
        } else {
          return Observable.of<SearchResponse>(null);
        }
      })
      .catch(error => this.onError(error))
      .subscribe((response: SearchResponse) => this.onSuccess(response));
  }

}
