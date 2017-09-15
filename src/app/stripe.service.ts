import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';

export interface PayResponse {
  error?: string
  scheduled: string[]
}

@Injectable()
export class StripeService {

  constructor(private http: HttpClient) { }

  pay(id: string, amount: number): Observable<PayResponse> {
    // TODO: XSRF protection.
    return this.http.post<PayResponse>('/api/donate', {id: id, amount: amount});
  }
}
