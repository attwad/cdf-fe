import { Component, OnInit } from '@angular/core';
import {Router, ActivatedRoute, ParamMap} from '@angular/router';

import {Lesson} from '../lesson';
import {LessonsService, LessonsResponse} from '../lessons.service';
import { ScrollerService } from '../scroller.service';

import 'rxjs/add/operator/switchMap';

@Component({
  selector: 'lessons',
  templateUrl: './lessons.component.html',
  styleUrls: ['./lessons.component.css']
})
export class LessonsComponent implements OnInit {
  lessons: Lesson[];
  cursor?: string;
  showConvertedOnly: boolean = false;

  constructor(
    private scrollerService: ScrollerService,
    private route: ActivatedRoute,
    private lessonsService: LessonsService,
    private router: Router) {}

  ngOnInit() {
    this.route.paramMap
      .switchMap((params: ParamMap) => {
        console.log('Params changed', params);
        this.showConvertedOnly = params.get('showConvertedOnly')==='true';
        return this.lessonsService.getLessons(
          this.showConvertedOnly, params.get('cursor'))
      }).subscribe((lessonsResponse: LessonsResponse) => {
        console.log("Fetched new response: ", lessonsResponse);
        this.lessons = lessonsResponse.lessons;
        this.cursor = lessonsResponse.cursor;
        console.log('fetched new lessons, cursor=', this.cursor);
        this.scrollerService.scrollToTop();
      });
  }

  onShowConvertedOnlyChanged(): void {
    this.router.navigate(
      ['/lessons', {
        showConvertedOnly: !this.showConvertedOnly
      }]);
  }

  goToNextPage():void {
    this.router.navigate(
      ['/lessons', {
        showConvertedOnly: this.showConvertedOnly,
        cursor: this.cursor
      }]);
  }
}
