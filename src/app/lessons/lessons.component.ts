import { Component, OnInit } from '@angular/core';
import {Router, ActivatedRoute, ParamMap} from '@angular/router';
import { Title } from '@angular/platform-browser';

import {Lesson} from '../lesson';
import {LessonsService, LessonsResponse} from '../lessons.service';
import { ScrollerService } from '../scroller.service';

import 'rxjs/add/operator/switchMap';

@Component({
  selector: 'app-lessons',
  templateUrl: './lessons.component.html',
  styleUrls: ['./lessons.component.css']
})
export class LessonsComponent implements OnInit {
  lessons: Lesson[];
  cursor?: string;
  showConvertedOnly = false;
  loading = false;

  constructor(
    title: Title,
    private scrollerService: ScrollerService,
    private route: ActivatedRoute,
    private lessonsService: LessonsService,
    private router: Router) {
      title.setTitle('Browse College de France audio transcripts');
    }

  ngOnInit() {
    this.route.paramMap
      .switchMap((params: ParamMap) => {
        console.log('Params changed', params);
        this.showConvertedOnly = params.get('showConvertedOnly') === 'true';
        this.loading = true;
        return this.lessonsService.getLessons(
          this.showConvertedOnly, params.get('cursor'));
      }).subscribe((lessonsResponse: LessonsResponse) => {
        console.log('Fetched new response: ', lessonsResponse);
        if (lessonsResponse.lessons.length === 0) {
          this.cursor = undefined;
        } else {
          this.lessons = lessonsResponse.lessons;
          this.cursor = lessonsResponse.cursor;
        }
        this.loading = false;
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

  goToNextPage(): void {
    this.router.navigate(
      ['/lessons', {
        showConvertedOnly: this.showConvertedOnly,
        cursor: this.cursor
      }]);
  }
}
