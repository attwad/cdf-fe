import { NgModule }             from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { LessonsComponent }   from './lessons/lessons.component';
import { SearchComponent }   from './search/search.component';

const routes: Routes = [
  { path: 'lessons', component: LessonsComponent },
  { path: 'search', component: SearchComponent },
  { path: '', redirectTo: '/lessons', pathMatch: 'full' }
];

@NgModule({
  imports: [ RouterModule.forRoot(routes, { enableTracing: false }) ],
  exports: [ RouterModule ]
})
export class AppRoutingModule {}
