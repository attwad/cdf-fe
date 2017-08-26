import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import {Router} from '@angular/router';

import { SearchComponent } from './search.component';

import { Observable } from 'rxjs/Observable';

import { LessonsService, SearchResponse} from '../lessons.service';

class MockLessonsService {
  search(query: string): Observable<SearchResponse> {
    return Observable.of<SearchResponse>();
 }
}

class RouterStub {
  navigateByUrl(url: string) { return url; }
}

describe('SearchComponent', () => {
  let component: SearchComponent;
  let fixture: ComponentFixture<SearchComponent>;
  const mockLessonsService: MockLessonsService = new MockLessonsService();

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SearchComponent ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [
        {provide: LessonsService, useValue: mockLessonsService},
        {provide: Router, useClass: RouterStub}
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SearchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
