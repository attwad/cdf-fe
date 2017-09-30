import { async, ComponentFixture, TestBed, inject } from '@angular/core/testing';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {Router, ActivatedRoute, ParamMap} from '@angular/router';
import {RouterTestingModule} from '@angular/router/testing';
import { Observable } from 'rxjs/Observable';

import { LessonsComponent } from './lessons.component';
import { LessonsService, LessonsResponse} from '../lessons.service';
import { ScrollerService } from '../scroller.service';
import { Lesson } from '../lesson';

import 'rxjs/add/observable/of';

class MockLessonsService {
  getLessons(showConvertedOnly: boolean, cursor?: string):
      Observable<LessonsResponse> {
    return Observable.of<LessonsResponse>({
      cursor: 'next cursor',
      lessons: []
    });
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
        { provide: ScrollerService, useValue: mockScrollerService}
      ]
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

  it('should have sane defaults', () => {
    expect(component.showConvertedOnly).toBeFalsy();
    expect(component.loading).toBeFalsy();
    expect(component.cursor).toBeUndefined();
  });

  it('should change route to show converted only', inject(
      [Router], (router: Router) => {
    let navigateSpy = spyOn(router, 'navigate');
    component.onShowConvertedOnlyChanged();
    expect(navigateSpy).toHaveBeenCalledWith(
      ['/lessons', {showConvertedOnly: true}]);
  }));

  it('should go to next page', inject([Router], (router: Router) => {
    let navigateSpy = spyOn(router, 'navigate');
    component.goToNextPage();
    // TODO: Wait for first async result to return, then cursor should be set.
    expect(navigateSpy).toHaveBeenCalledWith(
      ['/lessons', {showConvertedOnly: false, cursor: undefined}]);
  }));
});
