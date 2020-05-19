import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {LoginComponent} from "./components/login/login.component";
import {MyProfileComponent} from "./components/my-profile/my-profile.component";
import {RegisterComponent} from "./components/register/register.component";
import {SocietiesComponent} from "./components/societies/societies.component";
import {SocietyDetailsComponent} from "./components/society-details/society-details.component";
import {GoogleMapComponent} from "./components/google-map/google-map.component";
import {CreateTrashComponent} from "./components/create-trash/create-trash.component";
import {TrashDetailsComponent} from "./components/trash-details/trash-details.component";
import {EditTrashComponent} from "./components/edit-trash/edit-trash.component";
import {EditSocietyComponent} from "./components/edit-society/edit-society.component";
import {CreateEventComponent} from "./components/create-event/create-event.component";
import {EventDetailsComponent} from "./components/event-details/event-details.component";
import {EventsComponent} from "./components/events/events.component";

const routes: Routes = [
  { path: '', redirectTo:'map', pathMatch: 'full'},
  { path: 'report', component: CreateTrashComponent},
  { path: 'report/:lat/:lng', component: CreateTrashComponent},
  { path: 'trash/details/:id', component: TrashDetailsComponent},
  { path: 'trash/edit/:id', component: EditTrashComponent},
  { path: 'map', component: GoogleMapComponent},
  { path: 'me', component: MyProfileComponent},
  { path: 'societies', component: SocietiesComponent},
  {path: 'societies/:societyId', component: SocietyDetailsComponent},
  {path: 'societies/edit/:societyId', component: EditSocietyComponent},
  {path: 'events', component: EventsComponent},
  {path: 'events/create', component: CreateEventComponent},
  {path: 'events/details/:eventId', component: EventDetailsComponent},
  //{path: 'event/edit/:id', component}
  { path: 'register', component: RegisterComponent },
  { path: 'login', component: LoginComponent },
  //{ path: 'profile/:userId', component: UserDetailComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
