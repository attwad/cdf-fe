import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';

import { Observable }     from 'rxjs/Observable';

import {Lesson} from './lesson';

export interface LessonsResponse {
  cursor: string;
  lessons: Lesson[];
}

export interface Source {
  title: string;
	lecturer: string;
	chaire:     string
	type?:       string;
	type_title?: string;
	lang?:   string;
	source_url: string;
	transcript: string;
}

export interface SearchResponse {
  took_ms: number;
  timed_out: boolean;
  sources: Source[];
}

@Injectable()
export class LessonsService {

  constructor(private http: HttpClient) { }

  getLessons(showConvertedOnly: boolean, cursor?: string): Observable<LessonsResponse> {
    return this.http.get<LessonsResponse>(
      '/api/lessons', {
        params: new HttpParams()
            .set('cursor', cursor || '')
            .set('filter', showConvertedOnly ? 'converted': '')
      });
  }

  search(query: string): Observable<SearchResponse> {
    return this.http.get<SearchResponse>(
      'api/search', {params: new HttpParams().set('q', query)});
 }
}
