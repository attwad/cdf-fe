import { TestBed, inject } from '@angular/core/testing';
import {HttpClientTestingModule} from '@angular/common/http/testing';

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
});
