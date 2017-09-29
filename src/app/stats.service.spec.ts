import { TestBed, inject } from '@angular/core/testing';
import {HttpClientTestingModule} from '@angular/common/http/testing';

import { StatsService } from './stats.service';

describe('StatsService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [StatsService],
      imports: [HttpClientTestingModule]
    });
  });

  it('should be created', inject([StatsService], (service: StatsService) => {
    expect(service).toBeTruthy();
  }));
});
