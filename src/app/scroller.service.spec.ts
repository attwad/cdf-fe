import { TestBed, inject } from '@angular/core/testing';

import { ScrollerService } from './scroller.service';

describe('ScrollerService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ScrollerService]
    });
  });

  it('should be created', inject([ScrollerService], (service: ScrollerService) => {
    expect(service).toBeTruthy();
  }));
});
