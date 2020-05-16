import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {ApisModel} from "../../api/api-urls";
import {LoadImageProps, SetLoadedImageProps,} from 'ng-lazyload-image';
import {AuthService} from "../../services/auth/auth.service";

@Component({
  selector: 'app-image',
  template: `
    <img #element [defaultImage]="defaultImage" [lazyLoad]="image">
  `,
  styleUrls: ['./lazy-load-img.component.css']
})
export class LazyLoadImgComponent implements OnInit {
  @ViewChild('element', {static: false}) el: HTMLImageElement;

  url: string;

  @Input() image: string;
  @Input() service: string;

  defaultImage: string = 'https://cdn.onlinewebfonts.com/svg/img_258083.png';

  constructor(
    private authService: AuthService,
  ) {
  }

  ngOnInit() {
    if (this.service === ApisModel.user) {
      this.url = ApisModel.apiUrl + '/' + ApisModel.fileupload + '/' + ApisModel.user + '/' + this.image
    } else if (this.service === ApisModel.society){
      this.url = ApisModel.apiUrl + '/' + ApisModel.fileupload + '/' + ApisModel.society + '/' + this.image
    } else if (this.service === ApisModel.collection){
      this.url = ApisModel.apiUrl + '/' + ApisModel.fileupload + '/' + ApisModel.collection + '/' + this.image
    } else if (this.service === ApisModel.trash){
      this.url = ApisModel.apiUrl + '/' + ApisModel.fileupload + '/' + ApisModel.trash + '/' + this.image
    }
  }

  async ngAfterViewInit() {
    if (this.image === '' || this.image == null){
      return
    }

    if (this.el) {
      const imageBlob = await this.loadImage({
        element: this.el,
        imagePath: this.url
      })

      console.log('blob je: ', imageBlob)
      this.el.src = imageBlob
      this.setLoadedImage({element: this.el, imagePath: imageBlob})
    }
  }

  async loadImage({imagePath}: LoadImageProps) {
    const token = this.authService.getToken()
    return await fetch(imagePath, {
      headers: {
        Authorization: 'Bearer ' + token
      }
    }).then(res => res.blob()).then(blob => URL.createObjectURL(blob))
  }

  setLoadedImage({element, imagePath, useSrcset}: SetLoadedImageProps) {
    console.log(this.el.currentSrc)
    this.el.src = imagePath
    console.log(this.el.currentSrc)
  }

}
