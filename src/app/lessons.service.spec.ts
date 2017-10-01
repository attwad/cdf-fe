import { TestBed, inject, async } from '@angular/core/testing';
import {HttpClientTestingModule, HttpTestingController} from '@angular/common/http/testing';

import { LessonsService, LessonsResponse, SearchResponse } from './lessons.service';

describe('LessonsService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [LessonsService],
      imports: [HttpClientTestingModule]
    });
  });

  it('should be created', inject([LessonsService], (service: LessonsService) => {
    expect(service).toBeTruthy();
  }));

  it('should get converted lessons', async(inject([LessonsService, HttpTestingController],
      (service: LessonsService, httpMock: HttpTestingController) => {
    expect(service).toBeTruthy();
    service.getLessons(true, 'c1').subscribe((response: LessonsResponse) => {
      expect(response.cursor).toEqual('the next cursor');
    }, (error) => {fail('Error in getLessons:' + error.toString())}
    );
    const req = httpMock.expectOne('/api/lessons?cursor=c1&filter=converted');
    expect(req.request.method).toEqual('GET');
    req.flush({cursor: 'the next cursor'});
    httpMock.verify();
  })));

  it('should get non converted lessons', async(inject([LessonsService, HttpTestingController],
      (service: LessonsService, httpMock: HttpTestingController) => {
    expect(service).toBeTruthy();
    service.getLessons(false, 'c2').subscribe((response: LessonsResponse) => {
      expect(response.cursor).toEqual('the next cursor');
    }, (error) => {fail('Error in getLessons:' + error.toString())}
    );
    const req = httpMock.expectOne('/api/lessons?cursor=c2&filter=');
    expect(req.request.method).toEqual('GET');
    req.flush({cursor: 'the next cursor'});
    httpMock.verify();
  })));

  it('should search', async(inject([LessonsService, HttpTestingController],
      (service: LessonsService, httpMock: HttpTestingController) => {
    expect(service).toBeTruthy();
    service.search('a query', 5, 10).subscribe((response: SearchResponse) => {
      expect(response).toEqual({
        took_ms: 42,
        timed_out: false,
        total: 2,
        sources: []
      });
    }, (error) => {fail('Error in search:' + error.toString())}
    );
    const req = httpMock.expectOne('/api/search?q=a%20query&from=5&size=10');
    expect(req.request.method).toEqual('GET');
    req.flush({
      took_ms: 42,
      timed_out: false,
      total: 2,
      sources: []
    });
    httpMock.verify();
  })));
});
