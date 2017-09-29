import { TestBed, inject } from '@angular/core/testing';
import {HttpClientTestingModule} from '@angular/common/http/testing';

import { LessonsService } from './lessons.service';

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
});
