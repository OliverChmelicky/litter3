import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {SocietyService} from "../../services/society/society.service";
import {SocietyModel} from "../../models/society.model";

@Component({
  selector: 'app-society-details',
  templateUrl: './society-details.component.html',
  styleUrls: ['./society-details.component.css']
})
export class SocietyDetailsComponent implements OnInit {
  society: SocietyModel = {
    Id: '',
    Name: '',
    Avatar: '',
    Description: '',
    Users: [],
    CreatedAt: new Date()
  };

  constructor(
    private route: ActivatedRoute,
    private societyService: SocietyService,
  ) {
  }

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.societyService.getSociety(params.get('societyId')).subscribe(
        society => this.society = society
      )
    });
  }

}
