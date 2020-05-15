import {Component, Input, OnInit} from '@angular/core';
import {ApisModel} from "../../api/api-urls";
import { LoadImageProps } from 'ng-lazyload-image';
import {AuthService} from "../../services/auth/auth.service";
import {UserService} from "../../services/user/user.service";

@Component({
  selector: 'app-image',
  template: `
    <img [defaultImage]="defaultImage" [lazyLoad]="image">
  `,
  styleUrls: ['./lazy-load-img.component.css']
})
export class LazyLoadImgComponent implements OnInit {
  @Input() image;

  defaultImage: string = 'https://cdn.onlinewebfonts.com/svg/img_258083.png';

  constructor(
    private authService: AuthService,
  ) {
  }

  ngOnInit(): void {
    console.log('inputmam: ', this.image)
    this.image = ApisModel.pictureBucketPrefix + this.image
  }

  async loadImage({ imagePath }: LoadImageProps) {
    const token = this.authService.getToken()
    return await fetch(imagePath, {
      headers: {
        Authorization: 'Bearer ' + token
      }
    }).then(res => res.blob()).then(blob => URL.createObjectURL(blob));
  }

}
