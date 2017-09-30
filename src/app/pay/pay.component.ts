/// <reference types="stripe-checkout"/>
import { Component, OnInit } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { Title } from '@angular/platform-browser';

import { StripeService, PrepareResponse } from '../stripe.service';

interface proposition {
  hours: number
  priceUsdCents: number
}

@Component({
  selector: 'app-pay',
  templateUrl: './pay.component.html',
  styleUrls: ['./pay.component.css']
})
export class PayComponent implements OnInit {

  prepareResponse?: PrepareResponse;
  proposedHours = new Array<proposition>();
  loading = false;
  paid = false;
  step = 0;
  error: string;

  constructor(
    title: Title,
    private stripeService: StripeService) {
      title.setTitle('Contribute to College de France audio transcripts');
    }

  ngOnInit() {
    this.loading = true;
    this.stripeService.prepare()
    .subscribe((prepareResponse: PrepareResponse) => {
      console.log('Got prepare response', prepareResponse);
      this.loading = false;
      this.prepareResponse = prepareResponse;
      [1, 3, 5, 10].forEach((hours) => {
        this.proposedHours.push({
          hours: hours,
          priceUsdCents:
            (hours * prepareResponse.one_hour_amount_usd_cents)
            // Add stripe fees (3% + 30c).
            + (hours * prepareResponse.one_hour_amount_usd_cents * 0.03)
            + 30
        });
      })
      console.log("proposed payments:", this.proposedHours);
    },
    (err: HttpErrorResponse) => {
      console.error("error preparing:", err);
      this.error = err.message;
      this.loading = false;
    });
  }

  pay(hours, amount): void {
    const handler = (<any>window).StripeCheckout.configure({
      key: this.prepareResponse.stripe_publishable_key,
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
        },
        (err: HttpErrorResponse) => {
          console.error("error paying:", err);
          this.error = err.message;
          this.loading = false;
        });
      }
    }) as StripeCheckoutHandler;

    handler.open({
      name: 'college-audio.science',
      currency: 'USD',
      description: hours + 'H of audio transcriptions',
      amount: amount
    });
  }

  setStep(index: number) {
    this.step = index;
  }

}
