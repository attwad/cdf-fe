import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs/Observable';

export class Stats {
  NumConverted: number;
  NumTotal: number;
  ConvertedDurationSec: number;
  LeftDurationSec: number;
  Computed: string;
  PercentDone: number;
}

@Injectable()
export class StatsService {

  constructor(private http: HttpClient) { }

  getStats(): Observable<Stats> {
    return this.http.get<Stats>('/api/stats');
  }
}
