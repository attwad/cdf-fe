import { TestBed, inject } from '@angular/core/testing';
import { HttpClientModule } from '@angular/common/http';

import { LessonsService } from './lessons.service';

describe('LessonsService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [LessonsService],
      imports: [HttpClientModule]
    });
  });

  it('should be created', inject([LessonsService], (service: LessonsService) => {
    expect(service).toBeTruthy();
  }));
});
