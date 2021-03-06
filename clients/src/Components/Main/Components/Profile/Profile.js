import React, { Component } from 'react';
import api from '../../../../Constants/APIEndpoints/APIEndpoints';
import PageTypes from '../../../../Constants/PageTypes/PageTypes';
import Errors from '../../../Errors/Errors';
import Card from 'react-bootstrap/Card';
import Button from 'react-bootstrap/Button';
import ReactDOM from 'react-dom'



class Profile extends Component {
  constructor(props) {
      super(props);
      this.state = {
          userBadges: [],
          error: ''
      }
  }

  sendRequestThree = async (e) => {
    const response = await fetch(api.base + api.handlers.marketplaceBadges + this.props.user.id, {
      method: "GET",
      headers: new Headers({
        "Authorization": localStorage.getItem("Authorization")
      })
    });
    if (response.status >= 300) {
        const error = await response.text();
        this.setError(error);
        return;
    }
    const badgesResponse = await response.json();
    if (badgesResponse != null) {
      this.setState({
        userBadges: badgesResponse.map(badge => ({
          badgeID: badge.badgeID,
          cost: badge.cost,
          badgeName: badge.badgeName,
          badgeDescription: badge.badgeDescription,
          imgURL: badge.imgURL,
        }))
      });
    }   
  }

  componentWillMount() {
    {this.sendRequestThree()}
  }

  componentDidUpdate = () => { ReactDOM.findDOMNode(this).scrollIntoView(); }


  setError = (error) => {
      this.setState({ error })
  }

  render() {
      const { error} = this.state;
      const listItems = this.state.userBadges.map((badge) =>
        <li>
          <Card style={{ width: '15rem' }}>
            <Card.Img variant="top" src={badge.imgURL} />
            <Card.Body>
              <h2>{badge.badgeName}</h2>
            </Card.Body>
          </Card>
        </li>
      );
      return <div className="profile-page">
      <div className="profile">
          <Errors error={error} setError={this.setError} />
          <h1>Here's your <span className="red">Profile</span>:</h1>
          <Card style={{ width: '40rem' }}>
            <Card.Img variant="top" src={this.props.user.photoURL} />
            <Card.Body>
              <h2>{this.props.user.firstName} {this.props.user.lastName}</h2>
              <h4>Username: <span className="red">{this.props.user.userName}</span></h4>
              <Card.Text>
              <h4>Bio</h4>
                {this.props.user.bio}
              </Card.Text>
              <h4>Badges:</h4>
              <div className="badges">
                {listItems}
              </div>
              <Button variant="primary" onClick={(e) => { this.props.setPage(e, PageTypes.signedInUpdateName)}}>EDIT PROFILE</Button>
            </Card.Body>
          </Card>
          <div className="space" />
      </div>
      </div> 
  }

}

export default Profile;