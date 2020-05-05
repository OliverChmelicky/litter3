import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import {LoginComponent} from "./components/login/login.component";
import {MyProfileComponent} from "./components/my-profile/my-profile.component";
import {RegisterComponent} from "./components/register/register.component";
import {SocietiesComponent} from "./components/societies/societies.component";
import {SocietyDetailsComponent} from "./components/society-details/society-details.component";
import {GoogleMapComponent} from "./components/google-map/google-map.component";
import {CreateTrashComponent} from "./components/create-trash/create-trash.component";

const routes: Routes = [
  { path: '', redirectTo:'report', pathMatch: 'full'},
  { path: '**', component:  CreateTrashComponent},
  { path: 'report', component: CreateTrashComponent},
  { path: 'report/:lat/:lng', component: CreateTrashComponent},
  { path: 'map', component: GoogleMapComponent},
  { path: 'me', component: MyProfileComponent},
  { path: 'societies', component: SocietiesComponent},
  //{ path: 'profile/tu-daj-user-id', component: UserDetailComponent },
  {path: 'societies/:societyId', component: SocietyDetailsComponent},
  { path: 'register', component: RegisterComponent },
  { path: 'login', component: LoginComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
