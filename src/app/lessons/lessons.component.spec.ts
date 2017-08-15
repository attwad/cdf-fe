import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {Router, ActivatedRoute, ParamMap} from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { Observable } from 'rxjs/Observable';

import { LessonsComponent } from './lessons.component';
import { LessonsService, LessonsResponse} from '../lessons.service';

class MockLessonsService {
  getLessons(showConvertedOnly: boolean, cursor?: string): Observable<LessonsResponse> {
    return Observable.of<LessonsResponse>();
  }
}

describe('LessonsComponent', () => {
  let component: LessonsComponent;
  let fixture: ComponentFixture<LessonsComponent>;
  let mockLessonsService: MockLessonsService = new MockLessonsService();

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LessonsComponent ],
      imports: [RouterTestingModule.withRoutes(
        [{path: '', component: LessonsComponent}, {path: 'simple', component: LessonsComponent}]
      )],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [ {provide: LessonsService, useValue: mockLessonsService } ]
    })
    .compileComponents();
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
