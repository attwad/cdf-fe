import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import { PayComponent } from './pay.component';
import { Observable } from 'rxjs/Observable';
import {StripeService, PrepareResponse} from '../stripe.service';

class MockStripeService {
  pay(tokenId: string, email: string, amount: number): Observable<{}> {
    return Observable.of({});
  }

  prepare(): Observable<PrepareResponse> {
    return Observable.of({one_hour_amount_usd_cents: 200});
  }
}

describe('PayComponent', () => {
  let component: PayComponent;
  let fixture: ComponentFixture<PayComponent>;
  const mockStripeService: MockStripeService = new MockStripeService();

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PayComponent ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        {provide: StripeService, useValue: mockStripeService}
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PayComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('should have sane defaults', () => {
    expect(component.loading).toBeFalsy();
    expect(component.paid).toBeFalsy();
    expect(component.prepared).toBeTruthy();
  });

  it('Should get the initial cost per hour', () => {
    expect(component.oneHourAmountUsdCents).toEqual(200);
  });

  it('should set steps', () => {
    expect(component.step).toEqual(0);
    component.setStep(2);
    expect(component.step).toEqual(2);
  });
});
