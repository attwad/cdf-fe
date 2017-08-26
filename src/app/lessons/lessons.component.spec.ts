import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {Router, ActivatedRoute, ParamMap} from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { Observable } from 'rxjs/Observable';

import { LessonsComponent } from './lessons.component';
import { LessonsService, LessonsResponse} from '../lessons.service';
import { ScrollerService } from '../scroller.service';

class MockLessonsService {
  getLessons(showConvertedOnly: boolean, cursor?: string): Observable<LessonsResponse> {
    return Observable.of<LessonsResponse>();
  }
}

class MockScrollerService {
  public scrolled = false;

  scrollToTop(): void {
    this.scrolled = true;
  }
}

describe('LessonsComponent', () => {
  let component: LessonsComponent;
  let fixture: ComponentFixture<LessonsComponent>;
  const mockLessonsService: MockLessonsService = new MockLessonsService();
  const mockScrollerService: MockScrollerService = new MockScrollerService();

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LessonsComponent ],
      imports: [RouterTestingModule],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        { provide: LessonsService, useValue: mockLessonsService},
        { provide: ScrollerService, useValue: mockScrollerService}]
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LessonsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
