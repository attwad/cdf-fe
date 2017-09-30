import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';

export interface PrepareResponse {
  one_hour_amount_usd_cents: number;
  stripe_publishable_key: string;
}

@Injectable()
export class StripeService {

  constructor(private http: HttpClient) { }

  pay(tokenId: string, email: string, amountUsdCents: number): Observable<{}> {
    // Note: XSRF protection is handled at a higher level for us.
    // Cf. app.module.ts
    return this.http.post<{}>('/api/donate',
    {
      stripeToken: tokenId, stripeEmail: email,
      amountUsdCents: Math.floor(amountUsdCents)
    });
  }

  prepare(): Observable<PrepareResponse> {
    return this.http.get<PrepareResponse>('/api/donate');
  }
}
