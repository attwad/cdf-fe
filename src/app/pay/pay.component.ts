import { Component, OnInit } from '@angular/core';

import { StripeService, PayResponse } from '../stripe.service';

@Component({
  selector: 'app-pay',
  templateUrl: './pay.component.html',
  styleUrls: ['./pay.component.css']
})
export class PayComponent implements OnInit {

  payResponse?: PayResponse;

  constructor(private stripeService: StripeService) { }

  ngOnInit() {
  }

  pay(amount: number): void {
    let handler = (<any>window).StripeCheckout.configure({
      key: 'pk_test_VpXluQcCGGQVNgg0j0abLR5m',
      locale: 'auto',
      zipCode: true,
      allowRememberMe: false,
      token: (token: any) => {
        console.log('got token:', token);
        this.stripeService.pay(token.id, amount)
        .subscribe((payResponse: PayResponse) => {
          console.log('Got pay response:', payResponse);
          this.payResponse = payResponse;
        })
      }
    });

    handler.open({
      name: 'college-audio.science',
      currency: 'JPY',
      description: '~1H full lesson',
      amount: amount
    });
  }

}
