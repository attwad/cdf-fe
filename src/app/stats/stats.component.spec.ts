import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import {CUSTOM_ELEMENTS_SCHEMA} from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { By }          from '@angular/platform-browser';
import { DebugElement}  from '@angular/core';
import { StatsComponent } from './stats.component';
import { StatsService, Stats } from '../stats.service';

class MockStatsService {
  getStats(): Observable<Stats> {
    let stats = new Stats();
    stats.NumTotal = 5;
    stats.NumConverted = 2;
    stats.ConvertedDurationSec = 1000;
    stats.LeftDurationSec = 3000;
    return Observable.of<Stats>(stats);
  }
}

describe('StatsComponent', () => {
  let component: StatsComponent;
  let fixture: ComponentFixture<StatsComponent>;
  const mockStatsService: MockStatsService = new MockStatsService();
  let de:      DebugElement;
  let el:      HTMLElement;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StatsComponent ],
      schemas: [CUSTOM_ELEMENTS_SCHEMA],
      providers: [{ provide: StatsService, useValue: mockStatsService}]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StatsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });

  it('should output the right percentage', () => {
    de = fixture.debugElement.query(By.css('#percent'));
    el = de.nativeElement;
    expect(el.textContent).toBe('(25%)');
  });
});
