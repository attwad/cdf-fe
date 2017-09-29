/// <reference types="stripe-checkout"/>
import { Component, OnInit } from '@angular/core';

import { StripeService, PrepareResponse } from '../stripe.service';

@Component({
  selector: 'app-pay',
  templateUrl: './pay.component.html',
  styleUrls: ['./pay.component.css']
})
export class PayComponent implements OnInit {

  prepared = false;
  oneHourAmountUsdCents = 0;
  proposedHours = [1, 3, 5, 10];
  loading = false;
  paid = false;
  step = 0;

  constructor(private stripeService: StripeService) { }

  ngOnInit() {
    this.loading = true;
    this.stripeService.prepare()
    .subscribe((prepareResponse: PrepareResponse) => {
      console.log('Got prepare response', prepareResponse);
      this.loading = false;
      this.prepared = true;
      this.oneHourAmountUsdCents = prepareResponse.one_hour_amount_usd_cents;
    });
  }

  pay(numHours: number): void {
    const amount = numHours * this.oneHourAmountUsdCents;
    const handler = (<any>window).StripeCheckout.configure({
      key: 'pk_test_VpXluQcCGGQVNgg0j0abLR5m',
      locale: 'auto',
      zipCode: true,
      allowRememberMe: false,
      token: (token: any) => {
        this.loading = true;
        this.stripeService.pay(token.id, token.email, amount)
        .subscribe(() => {
          console.log('Got pay response');
          this.paid = true;
          this.loading = false;
        });
      }
    }) as StripeCheckoutHandler;

    handler.open({
      name: 'college-audio.science',
      currency: 'USD',
      description: numHours + 'H of audio transcriptions',
      amount: amount
    });
  }

  setStep(index: number) {
    this.step = index;
  }

}
