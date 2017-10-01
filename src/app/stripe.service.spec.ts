import { TestBed, inject } from '@angular/core/testing';
import {HttpClientTestingModule, HttpTestingController} from '@angular/common/http/testing';

import { StripeService, PrepareResponse } from './stripe.service';

describe('StripeService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [StripeService],
      imports: [HttpClientTestingModule]
    });
  });

  it('should be created', inject([StripeService], (service: StripeService) => {
    expect(service).toBeTruthy();
  }));

  it('Should send payments', inject([StripeService, HttpTestingController], (service: StripeService, httpMock: HttpTestingController) => {
    service.pay('token', 'foo@bar.com', 123.45).subscribe(() => {
    }, (err) => fail('Failed payment:' + err));
    const req = httpMock.expectOne('/api/donate');
    expect(req.request.method).toEqual('POST');
    expect(req.request.body).toEqual({
      stripeToken: 'token',
      stripeEmail: 'foo@bar.com',
      amountUsdCents: 123
    });
    req.flush({});
    httpMock.verify();
  }));

  it('Should prepare payments', inject([StripeService, HttpTestingController], (service: StripeService, httpMock: HttpTestingController) => {
    service.prepare().subscribe((resp: PrepareResponse) => {
      expect(resp.one_hour_amount_usd_cents).toEqual(200);
    }, (err) => fail('Failed payment preparation:' + err));
    const req = httpMock.expectOne('/api/donate');
    expect(req.request.method).toEqual('GET');
    req.flush({one_hour_amount_usd_cents: 200});
    httpMock.verify();
  }));
});
