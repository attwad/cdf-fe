import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { LessonsComponent } from './lessons/lessons.component';
import { SearchComponent } from './search/search.component';
import { AboutComponent } from './about/about.component';
import { PayComponent } from './pay/pay.component';

const routes: Routes = [
  { path: 'lessons', component: LessonsComponent },
  { path: 'search', component: SearchComponent },
  { path: 'about', component: AboutComponent },
  { path: 'donate', component: PayComponent },
  { path: '', redirectTo: '/search', pathMatch: 'full' }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes, { enableTracing: false }) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule {}
