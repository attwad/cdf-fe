import { TestBed, inject } from '@angular/core/testing';
import {HttpClientTestingModule, HttpTestingController} from '@angular/common/http/testing';

import { StatsService, Stats } from './stats.service';

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

  it('should fetch the stats', inject([StatsService, HttpTestingController], (service: StatsService, httpMock: HttpTestingController) => {
    service.getStats().subscribe((stats: Stats) => {
      expect(stats.NumConverted).toEqual(42);
    });
    const req = httpMock.expectOne('/api/stats');
    expect(req.request.method).toEqual('GET');
    req.flush({
      NumConverted: 42
    });
    httpMock.verify();
  }));
});
